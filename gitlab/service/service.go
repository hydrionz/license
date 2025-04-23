package service

import (
	"archive/zip"
	"crypto/aes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	gorsa "github.com/Lyafei/go-rsa"
	"github.com/gin-gonic/gin"
	"io"
	"license/config"
	"license/crypto"
	"license/gitlab/entity"
	"log"
	"net/http"
	"os"
	"time"
)

var privateKey string
var publicKey string

// LoadKeys reads, decodes and parses RSA private and public keys.
func LoadKeys() {
	// Read public key
	publicBytes, err := os.ReadFile(config.GetConfig().DataDir + "/.license_encryption_key.pub")
	if err != nil {
		log.Printf("Failed to read public key file: %v", err)
		return
	}
	// Convert to string
	publicKey = string(publicBytes)

	// Read private key
	privateBytes, err := os.ReadFile(config.GetConfig().DataDir + "/.license_decryption_key.pri")
	if err != nil {
		log.Printf("Failed to read private key file: %v", err)
		return
	}
	// Convert to string
	privateKey = string(privateBytes)
}

// createLicenseJson creates a JSON representation of the license
func createLicenseJson(licenseInfo entity.LicenseInfo, expireTime string) (string, error) {

	var expirationDate time.Time
	var err error
	if len(expireTime) == 0 {
		// Default expiration time is 2 years
		expirationDate = time.Now().AddDate(2, 0, 0)
	} else {
		expirationDate, err = time.Parse(time.DateTime, expireTime)
		if err != nil {
			log.Printf("Failed to parse expiration time: %v", err)
			return "", err
		}
	}

	license := entity.License{
		Version:                      1,
		License:                      licenseInfo,
		StartsAt:                     entity.CustomTime{Time: time.Now()},
		ExpiresAt:                    entity.CustomTime{Time: expirationDate},
		NotifyAdminsAt:               entity.CustomTime{Time: expirationDate},
		NotifyUsersAt:                entity.CustomTime{Time: expirationDate},
		BlockChangesAt:               entity.CustomTime{Time: expirationDate},
		CloudLicensingEnabled:        false,
		OfflineCloudLicensingEnabled: false,
		AutoRenewEnabled:             false,
		SeatReconciliationEnabled:    false,
		OperationalMetricsEnabled:    false,
		GeneratedFromCustomersDot:    false,
		Restrictions: entity.Restriction{
			ActiveUserCount: 10000,
			Plan:            "ultimate",
		},
	}

	jsonData, err := json.Marshal(license)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// generateRandomIV generates a random initialization vector (IV)
func generateRandomIV() ([]byte, error) {
	iv := make([]byte, aes.BlockSize) // AES block size is fixed at 16 bytes
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil
}

// Encrypt wraps the Encrypt method, using AES-CBC encryption and PKCS7 padding
func Encrypt(data, key, iv []byte) ([]byte, error) {
	aesTool := crypto.AesCbcPkcs7{Key: key, Iv: iv}
	enc, err := aesTool.Encrypt(data)
	if err != nil {
		log.Println("Encrypt error:", err)
		return nil, err
	}
	return enc, err
}

// Uses RSA private key to "encrypt" data
func encryptWithPrivateKey(data string) (string, error) {
	encrypt, err := gorsa.PriKeyEncrypt(data, privateKey)
	if err != nil {
		log.Printf("Failed to encrypt data with RSA private key: %v", err)
		return "", err
	}
	return encrypt, nil
}

// encryptLicense encrypts license data using AES and RSA
func encryptLicense(data string) (string, error) {
	// Generate 256-bit AES key
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		log.Printf("Failed to generate AES key: %v", err)
		return "", err
	}

	// Generate random IV
	iv, err := generateRandomIV()
	if err != nil {
		log.Printf("Failed to generate AES IV: %v", err)
		return "", err
	}

	encryptedData, err := Encrypt([]byte(data), key, iv)
	if err != nil {
		log.Printf("Failed to encrypt data: %v", err)
		return "", err
	}

	// Note: RSA encryption is typically done with a public key, but technically can be done with a private key (although not recommended)
	encryptedKey, err := encryptWithPrivateKey(string(key))
	if err != nil {
		return "", err
	}

	// Encode encrypted data as Base64
	encryptedDataStr := base64.StdEncoding.EncodeToString(encryptedData)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	// Package as JSON format
	result := map[string]string{
		"data": encryptedDataStr,
		"key":  encryptedKey,
		"iv":   ivStr,
	}
	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Printf("Failed to package JSON data: %v", err)
		return "", err
	}

	// Encode JSON as Base64
	encodedFinal := base64.StdEncoding.EncodeToString(jsonData)
	return encodedFinal, nil
}

// Generate generates a license and sends it via HTTP response
func Generate(ctx *gin.Context, licenseInfo entity.LicenseInfo, expireTime string) {
	createLicense(ctx, licenseInfo, expireTime)
}

// createLicense creates and sends a license
func createLicense(ctx *gin.Context, licenseInfo entity.LicenseInfo, expireTime string) {
	// Create license JSON data
	licenseJson, err := createLicenseJson(licenseInfo, expireTime)
	if err != nil {
		log.Printf("Failed to create license JSON: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Encrypt the license data
	encryptedLicense, err := encryptLicense(licenseJson)
	if err != nil {
		log.Printf("Failed to encrypt license: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Export a ZIP file containing the encrypted license and public key file
	err = exportZipStream(ctx, encryptedLicense)
	if err != nil {
		log.Printf("Failed to export ZIP file: %v", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
}

// exportZipStream creates and sends a ZIP file containing the encrypted license and public key file
func exportZipStream(ctx *gin.Context, encryptedLicense string) error {
	// Set response headers for file download
	ctx.Status(http.StatusOK) // Explicitly set status code to 200 OK
	ctx.Header("Content-Disposition", "attachment; filename=license.zip")
	ctx.Header("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(ctx.Writer)
	defer func(zipWriter *zip.Writer) {
		err := zipWriter.Close()
		if err != nil {
			log.Printf("Failed to close ZIP writer: %v", err)
		}
	}(zipWriter)

	// Add public key file to ZIP
	if err := addFileToZip(zipWriter, config.GetConfig().DataDir+"/.license_encryption_key.pub", "license/.license_encryption_key.pub"); err != nil {
		return err
	}

	// Add encrypted license data to ZIP
	if err := addLicenseToZip(zipWriter, encryptedLicense, "license/license.gitlab-license"); err != nil {
		return err
	}

	return nil
}

// addFileToZip reads a file from the filesystem and adds it to the ZIP
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(fileToZip *os.File) {
		err := fileToZip.Close()
		if err != nil {
			log.Printf("Failed to close file: %v", err)
		}
	}(fileToZip)

	// Get file info for setting the ZIP entry's size and timestamp
	fileInfo, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	// Create ZIP entry and manually set file timestamp
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = zipPath
	// Set compression method
	header.Method = zip.Deflate
	// Preserve original file's modification time
	header.Modified = fileInfo.ModTime()

	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write file data to ZIP
	_, err = io.Copy(zipFile, fileToZip)
	return err
}

// addLicenseToZip directly writes string data to a ZIP entry
func addLicenseToZip(zipWriter *zip.Writer, data, zipPath string) error {
	// Create a new zip.FileHeader, set filename and modification time
	header := &zip.FileHeader{
		Name:     zipPath,
		Method:   zip.Deflate, // Use compression to reduce file size
		Modified: time.Now(),  // Set current time as file modification time
	}

	// Create ZIP entry
	zipFile, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Write data to ZIP entry
	_, err = zipFile.Write([]byte(data))
	return err
}
