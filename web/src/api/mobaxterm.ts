import api from './config';
import axios from 'axios';

/**
 * Fetch available MobaXterm versions from the server
 * @returns Array of version strings
 */
export const fetchVersions = async (queryParam?: string): Promise<string[]> => {
  try {
    console.log('Calling /mobaxterm/versions API');
    const url = `/mobaxterm/versions${queryParam || ''}`;
    const response = await api.get(url, {
      headers: {
        'Cache-Control': 'no-cache, no-store',
        'Pragma': 'no-cache',
        'Expires': '0'
      }
    });
    console.log('API response:', response);

    // Backend directly returns an array
    if (Array.isArray(response)) {
      console.log('Got array response directly:', response);
      return response;
    }

    console.warn('Unexpected response format:', response);
    return [];
  } catch (error) {
    console.error('Error fetching versions:', error);
    return [];
  }
};

/**
 * Fetch available MobaXterm versions from the server
 * @returns Array of version strings
 */
export const generateLicense = async (
    name: string,
    version: string,
    count: string
): Promise<boolean> => {
  try {
    // Create FormData
    const formData = new FormData();
    formData.append('name', name);
    formData.append('version', version);
    formData.append('count', count);

    // 使用axios直接请求，以便可以设置responseType为blob
    const response = await axios({
      baseURL: "/api",
      method: 'post',
      url: '/mobaxterm/generate',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      responseType: 'blob' // 指定响应类型为blob
    });

    // 验证响应是否包含有效的blob数据
    if (response.data && response.data.size > 0) {
      // 创建Blob并生成下载链接
      const blob = new Blob([response.data], {
        type: response.headers['content-type'] || 'application/octet-stream'
      });
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;

      // 从Content-Disposition获取文件名，如果没有则使用默认名称
      const contentDisposition = response.headers['content-disposition'];
      let filename = 'Custom.mxtpro';
      if (contentDisposition) {
        const filenameMatch = contentDisposition.match(/filename="?([^"]+)"?/);
        if (filenameMatch && filenameMatch[1]) {
          filename = filenameMatch[1];
        }
      }

      console.log(filename);
      link.setAttribute('download', filename);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);

      return true;
    } else {
      throw new Error('返回的文件数据无效');
    }
  } catch (error) {
    console.error('生成许可证失败:', error);
    throw error;
  }
};