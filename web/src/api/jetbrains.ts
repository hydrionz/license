import api from './config';

export interface JetBrainsGenerateResponse {
  activationCode: string;
  powerConfig: string;
  licenseId: string;
  expiresAt: string;
  generatedAt: string;
}

/**
 * 生成JetBrains激活码
 * @param licenseeName 授权用户名
 * @param effectiveDate 有效日期，可选，格式为: YYYY-MM-DD HH:MM:SS
 * @param codes 产品代码，多个代码用逗号分隔
 */
export const generateLicense = async (
  licenseeName?: string, 
  effectiveDate?: string,
  codes?: string
): Promise<JetBrainsGenerateResponse> => {
  const params: any = {};
  
  if (licenseeName) params.licenseeName = licenseeName;
  if (effectiveDate) params.effectiveDate = effectiveDate;
  if (codes) params.codes = codes;
  
  return api.get<JetBrainsGenerateResponse>('/jetbrains/generate', {params: params});
};

/**
 * 获取JetBrains服务器规则
 */
export const getLicenseServerRule = async (): Promise<string> => {
  return api.getRaw<string>('/jetbrains/licenseServerRule');
};
