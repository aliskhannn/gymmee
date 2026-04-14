import axios from 'axios';
import WebApp from '@twa-dev/sdk';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

export const apiClient = axios.create({
  baseURL: API_URL,
});

apiClient.interceptors.request.use((config) => {
  const initData = WebApp.initData;
  
  if (initData) {
    config.headers.Authorization = `tma ${initData}`;
  }
  
  return config;
}, (error) => {
  return Promise.reject(error);
});