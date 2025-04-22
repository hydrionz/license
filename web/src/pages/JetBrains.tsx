import React, { useState, useEffect } from 'react';
import { Typography, Form, Button, Alert, Card, Space, message, Input, Radio } from 'antd';
import styled from 'styled-components';
import { 
  LoadingOutlined, 
  InfoCircleOutlined,
  CopyOutlined,
  CheckOutlined
} from '@ant-design/icons';
import PageHeader from '../components/PageHeader';
import ResultCard from '../components/ResultCard';
import { jetbrains } from '../api';
import { JetBrainsLicense } from '../types';

const { Paragraph, Text, Title } = Typography;

const FormCard = styled(Card)`
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 32px;
  border: 1px solid #e5e7eb;
  
  .ant-card-head {
    border-bottom: 1px solid #e5e7eb;
  }
`;

const SubmitButton = styled(Button)`
  width: 100%;
  height: 40px;
  border-radius: 8px;
  margin-top: 8px;
`;

const LicenseContent = styled.div`
  margin-top: 12px;
  background-color: #f9fafb;
  padding: 16px;
  padding-right: 48px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  white-space: pre-wrap;
  word-break: break-all;
  position: relative;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 14px;
  overflow-x: auto;
`;

const CopyButton = styled(Button)`
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0.8;
  z-index: 2;
  
  &:hover {
    opacity: 1;
  }
`;

const LicenseResultCard = styled(Card)`
  margin-bottom: 32px;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid #e5e7eb;
  overflow: hidden;

  .ant-card-head {
    background-color: #f9fafb;
    border-bottom: 1px solid #e5e7eb;
  }
`;

const LabelText = styled(Text)`
  font-weight: 500;
  color: #4b5563;
`;

const JetBrains: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [license, setLicense] = useState<JetBrainsLicense | null>(null);
  const [rawResponse, setRawResponse] = useState<string | null>(null);
  const [serverRule, setServerRule] = useState<string>('');
  const [loadingServerRule, setLoadingServerRule] = useState(false);
  const [activationMethod, setActivationMethod] = useState<'code' | 'server'>('code');
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [codeForm] = Form.useForm();

  // 当用户切换到服务器激活模式时，获取服务器规则
  useEffect(() => {
    if (activationMethod === 'server' && !serverRule && !loadingServerRule) {
      setLoadingServerRule(true);
      
      const fetchServerRule = async () => {
        try {
          const serverRuleText = await jetbrains.getLicenseServerRule();
          setServerRule(serverRuleText);
        } catch (error) {
          console.error('获取服务器规则失败:', error);
        } finally {
          setLoadingServerRule(false);
        }
      };
  
      fetchServerRule();
    }
  }, [activationMethod, serverRule, loadingServerRule]);

  const handleGenerateLicense = async (values: { 
    licenseeName?: string, 
    effectiveDate?: string, 
    codes?: string 
  }) => {
    setLoading(true);
    try {
      const data = await jetbrains.generateLicense(
        values.licenseeName, 
        values.effectiveDate, 
        values.codes
      );

      // Store the raw response data for our custom rendering
      setRawResponse(typeof data === 'string' ? data : String(data));
      
      // Also keep the license object for compatibility
      setLicense({
        code: '',
        product: values.codes?.split(',')[0] || 'Unknown'
      });
    } catch (error) {
      console.error('生成许可证失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = (key: string, text: string) => {
    navigator.clipboard.writeText(text).then(() => {
      setCopying({ ...copying, [key]: true });
      message.success('复制成功');
      
      setTimeout(() => {
        setCopying({ ...copying, [key]: false });
      }, 2000);
    }).catch(() => {
      message.error('复制失败，请手动复制');
    });
  };

  // Extract power.conf content from raw response
  const extractPowerConf = (): string | null => {
    if (!rawResponse) return null;
    
    // Find the power.conf section
    const startMarker = "================== power.conf ==================";
    const endMarker = "================== power.conf ==================";
    
    const startIdx = rawResponse.indexOf(startMarker);
    if (startIdx === -1) return null;
    
    const contentStart = startIdx + startMarker.length;
    const endIdx = rawResponse.indexOf(endMarker, contentStart);
    if (endIdx === -1) return null;
    
    return rawResponse.substring(contentStart, endIdx).trim();
  };

  // Extract activation code from raw response
  const extractActivationCode = (): string | null => {
    if (!rawResponse) return null;
    
    // Find the activation code section
    const startMarker = "================== activation code ==================";
    const endMarker = "================== activation code ==================";
    
    const startIdx = rawResponse.indexOf(startMarker);
    if (startIdx === -1) return null;
    
    const contentStart = startIdx + startMarker.length;
    const endIdx = rawResponse.indexOf(endMarker, contentStart);
    if (endIdx === -1) return null;
    
    return rawResponse.substring(contentStart, endIdx).trim();
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: '首页',
    },
    {
      path: '',
      breadcrumbName: 'JetBrains 激活工具',
    },
  ];

  const onFinish = (values: any) => {
    if (!values.licenseeName) {
      message.error('请输入授权用户名');
      return;
    }

    // 如果是按照产品生成激活码
    if (activationMethod === 'code') {
      if (!values.manualCodes) {
        message.error('请输入产品代码');
        return;
      }
      
      handleGenerateLicense({
        licenseeName: values.licenseeName,
        effectiveDate: values.effectiveDate,
        codes: values.manualCodes
      });
    } else if (activationMethod === 'server') {
      // 服务器激活方式，不需要传递产品代码
      handleGenerateLicense({
        licenseeName: values.licenseeName,
        effectiveDate: values.effectiveDate
      });
    }
  };

  return (
    <div>
      <PageHeader
        title="JetBrains 激活工具"
        subTitle="生成JetBrains全系列产品的激活码"
        breadcrumbs={breadcrumbs}
      />

      <FormCard title="JetBrains 激活工具">
        <Radio.Group 
          value={activationMethod} 
          onChange={(e) => setActivationMethod(e.target.value)}
          style={{ marginBottom: 16 }}
        >
          <Radio.Button value="code">激活码激活</Radio.Button>
          <Radio.Button value="server">在线服务器激活</Radio.Button>
        </Radio.Group>

        {activationMethod === 'code' ? (
          <Form form={codeForm} onFinish={onFinish} layout="vertical">
            <Form.Item
              name="licenseeName"
              label="授权用户名"
              rules={[{ required: true, message: '请输入授权用户名' }]}
            >
              <Input placeholder="请输入授权用户名" />
            </Form.Item>

            <Form.Item
              name="effectiveDate"
              label="有效日期"
            >
              <Input placeholder="例如: 2024-05-01 12:30:00" />
            </Form.Item>
            
            <Form.Item
              name="manualCodes"
              label="产品代码"
              rules={[{ required: true, message: '请输入产品代码' }]}
            >
              <Input.TextArea 
                placeholder="请输入产品代码，多个产品用逗号分隔" 
                rows={3}
              />
            </Form.Item>

            <Form.Item>
              <SubmitButton
                type="primary"
                htmlType="submit"
                loading={loading}
              >
                生成激活码
              </SubmitButton>
            </Form.Item>
          </Form>
        ) : (
          <div>
            <Paragraph>
              您也可以通过配置激活服务器的方式激活JetBrains产品。复制下面的power.conf配置到JetBrains激活服务器设置中：
            </Paragraph>

            {serverRule ? (
              <ResultCard
                title="服务器激活配置"
                data={{
                  'power.conf配置': serverRule,
                }}
                fileName="jetbrains-server-config.txt"
              />
            ) : (
              <Alert
                message={loadingServerRule ? "正在加载服务器规则，请稍候..." : "点击服务器激活选项后将自动加载服务器规则"}
                type="info"
                showIcon
                icon={loadingServerRule ? <LoadingOutlined /> : <InfoCircleOutlined />}
              />
            )}
          </div>
        )}
      </FormCard>

      {/* Custom license result display */}
      {license && rawResponse && (
        <LicenseResultCard title={<Title level={5} style={{ margin: 0 }}>激活码生成成功</Title>}>
          <Space direction="vertical" style={{ width: '100%' }}>
            <div>
              <LabelText>产品:</LabelText>
              <LicenseContent>
                {license.product || '未知产品'}
              </LicenseContent>
            </div>
            
            {extractPowerConf() && (
              <div style={{ marginTop: 16 }}>
                <LabelText>power.conf配置:</LabelText>
                <LicenseContent>
                  {extractPowerConf()}
                  <CopyButton
                    size="small"
                    type="primary"
                    ghost
                    icon={copying['powerConf'] ? <CheckOutlined /> : <CopyOutlined />}
                    onClick={() => copyToClipboard('powerConf', extractPowerConf() || '')}
                  />
                </LicenseContent>
              </div>
            )}
            
            {extractActivationCode() && (
              <div style={{ marginTop: 16 }}>
                <LabelText>激活码:</LabelText>
                <LicenseContent>
                  {extractActivationCode()}
                  <CopyButton
                    size="small"
                    type="primary"
                    ghost
                    icon={copying['activationCode'] ? <CheckOutlined /> : <CopyOutlined />}
                    onClick={() => copyToClipboard('activationCode', extractActivationCode() || '')}
                  />
                </LicenseContent>
              </div>
            )}
          </Space>
        </LicenseResultCard>
      )}
    </div>
  );
};

export default JetBrains; 