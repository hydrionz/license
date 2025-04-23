import api from './config';
import { MobaXtermLicense } from '../types';

/**
 * Generate MobaXterm license
 * @param formData Form data containing username, version, count
 * @returns File stream for the license
 */
export const generateLicense = async (
  formData: FormData
): Promise<Blob> => {
  return api.post('/mobaxterm/generate', formData, {
    responseType: 'blob',
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  });
};

/**
 * Fetch available MobaXterm versions from the server
 * @returns Array of version strings
 */
export const fetchVersions = async (): Promise<string[]> => {
  try {
    console.log('Calling /mobaxterm/versions API');
    const response = await api.get('/mobaxterm/versions');
    console.log('API response:', response);
    
    // Extract data from the response
    if (response && response.data) {
      if (Array.isArray(response.data)) {
        // Direct array response
        console.log('Direct array response:', response.data);
        return response.data;
      } else if (response.data.data && Array.isArray(response.data.data)) {
        // Nested data object
        console.log('Nested data object:', response.data.data);
        return response.data.data;
      } else if (typeof response.data === 'object') {
        // Try to extract any array in the response
        console.log('Looking for arrays in response');
        for (const key in response.data) {
          if (Array.isArray(response.data[key])) {
            console.log(`Found array in key ${key}:`, response.data[key]);
            return response.data[key];
          }
        }
      }
    }
    
    console.warn('Could not extract versions from response', response);
    return [];
  } catch (error) {
    console.error('Error fetching versions:', error);
    return [];
  }
}; 