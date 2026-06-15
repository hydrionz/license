import axios, {AxiosRequestConfig} from 'axios';
import {message} from 'antd';
import {ApiResponse} from '../types';

// 创建 axios 实例
const axiosInstance = axios.create({
  baseURL: '/api', // 在实际部署时，可以通过环境变量或配置文件设置
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

const getErrorMessage = (error: any): string => {
  return error.response?.data?.message || error.message || '网络错误，请稍后重试';
};

const unwrapResponse = <T>(responseData: ApiResponse<T>): T => {
  if (responseData.code !== 200) {
    throw new Error(responseData.message || '请求失败');
  }
  return responseData.data;
};

// 提供一个通用的API请求函数
const api = {
  async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axiosInstance.get(url, config);
      return unwrapResponse(response.data as ApiResponse<T>);
    } catch (error: any) {
      message.error(getErrorMessage(error));
      throw error;
    }
  },

  async post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axiosInstance.post(url, data, config);
      return unwrapResponse(response.data as ApiResponse<T>);
    } catch (error: any) {
      message.error(getErrorMessage(error));
      throw error;
    }
  },

  // 原始响应：用于非标准接口（例如纯文本配置）
  async getRaw<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axiosInstance.get(url, config);
      return response.data as T;
    } catch (error: any) {
      message.error(getErrorMessage(error));
      throw error;
    }
  },

  async postRaw<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> {
    try {
      const response = await axiosInstance.post(url, data, config);
      return response.data as T;
    } catch (error: any) {
      message.error(getErrorMessage(error));
      throw error;
    }
  }
};

export default api; 
