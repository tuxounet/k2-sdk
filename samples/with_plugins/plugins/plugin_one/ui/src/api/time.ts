import axios from 'axios';
import { TimeResponse } from '../types/time';

const API_BASE = (import.meta as any).env.DEV ? 'http://localhost:8080/horloge/temps' : '/horloge/temps';

export const getTime = async (): Promise<TimeResponse> => {
  const response = await axios.get<TimeResponse>(`${API_BASE}/time`);
  return response.data;
};
