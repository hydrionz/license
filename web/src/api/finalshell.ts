import api from "./config";

export interface FinalShellLicense {
    advancedBelow396: string;
    proBelow396: string;
    advancedAbove396: string;
    proAbove396: string;
    advancedAbove45: string;
    proAbove45: string;
    advancedAbove46: string;
    proAbove46: string;
}

/**
 * Generate FinalShell license
 */
export const generateLicense = async (machineCode: string): Promise<FinalShellLicense> => {
  const formData = new FormData();
  formData.append('machineCode', machineCode);
 return  api.post<FinalShellLicense>('/final-shell/generateLicense', formData, {
     headers: {
         'Content-Type': 'multipart/form-data'
     }
 })
}; 