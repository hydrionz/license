import api from './config';
import axios from 'axios';

/**
 * 生成FinalShell许可证
 * @param machineCode 机器码
 */
export const generateLicense = async (machineCode: string): Promise<string[]> => {
  const formData = new FormData();
  formData.append('machineCode', machineCode);
  
  // 直接使用 axios 发送请求，确保正确传递 FormData
  const response = await axios.post('/api/final-shell/generateLicense', formData);
  return response.data;
}; 