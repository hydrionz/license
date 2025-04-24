import api from './config';

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