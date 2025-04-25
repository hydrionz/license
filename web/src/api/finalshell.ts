import api from "./config";

/**
 * Generate FinalShell license
 * @param machineCode Machine code
 * @returns An array of license codes in order: 
 *   [0] - Version < 3.9.6 Advanced Edition 
 *   [1] - Version < 3.9.6 Professional Edition
 *   [2] - Version >= 3.9.6 Advanced Edition
 *   [3] - Version >= 3.9.6 Professional Edition
 */
export const generateLicense = async (machineCode: string): Promise<string[]> => {
  const formData = new FormData();
  formData.append('machineCode', machineCode);
 return  api.post<string[]>('/final-shell/generateLicense', formData, {
     headers: {
         'Content-Type': 'multipart/form-data'
     }
 })
}; 