import React, { useState } from 'react';
import {Typography, Form, Button, Input, Card, App, Row, Col} from 'antd';
import styled from 'styled-components';
import { useTranslation } from 'react-i18next';
import { CopyOutlined, CheckOutlined } from '@ant-design/icons';
import PageHeader from '../components/PageHeader';
import UsageNotice from '../components/UsageNotice';
import { finalshell } from '../api';
import { copyAndManageState } from '../utils/clipboard';

const { Paragraph } = Typography;

const FormWrapper = styled.div`
  width: 100%;
  margin-bottom: 32px;
`;

const FormCard = styled(Card)`
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 32px;
  border: 1px solid #e5e7eb;
  
  .ant-card-head {
    border-bottom: 1px solid #e5e7eb;
  }
`;

const StepItem = styled.div`
  margin-bottom: 16px;
  display: flex;
  align-items: flex-start;
`;

const StepNumber = styled.span`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  background-color: #1890ff;
  color: #fff;
  border-radius: 50%;
  margin-right: 12px;
  font-size: 14px;
  flex-shrink: 0;
`;

const StepContent = styled.div`
  flex: 1;
`;

const AuthorizationCodeContainer = styled.div`
  position: relative;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
  padding-right: 50px;
  margin-bottom: 0;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  overflow-wrap: break-word;
  word-break: break-all;
  
  @media (max-width: 768px) {
    padding: 12px;
    padding-right: 40px;
    font-size: 12px;
  }
`;

const ButtonContainer = styled.div`
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  
  @media (max-width: 768px) {
    top: 4px;
    right: 4px;
  }
`;

const CopyButton = styled(Button)`
  opacity: 0.8;
  
  &:hover {
    opacity: 1;
  }
  
  @media (max-width: 768px) {
    padding: 0 8px;
    height: 24px;
    font-size: 12px;
  }
`;

const CodeContainer = styled.div`
  position: relative;
  margin-bottom: 16px;
  overflow: hidden;
  isolation: isolate;
`;

const CodeLabel = styled.div`
  font-weight: 600;
  margin-bottom: 4px;
  color: #4b5563;
`;

const FinalShell: React.FC = () => {
  const { t } = useTranslation();
  const { notification } = App.useApp();
  const [loading, setLoading] = useState(false);
  const [authorizationCodes, setAuthorizationCodes] = useState<string[]>([]);
  const [copying, setCopying] = useState<{[key: string]: boolean}>({});
  const [form] = Form.useForm();

  // Version labels in order matching the backend response
  const versionLabels = [
    'finalshell.versions.advancedBelow396',
    'finalshell.versions.proBelow396',
    'finalshell.versions.advancedAbove396',
    'finalshell.versions.proAbove396',
    'finalshell.versions.advancedAbove45',
    'finalshell.versions.proAbove45'
  ];

  const handleGenerateLicense = async (values: { machineCode: string }) => {
    setLoading(true);
    try {
      const data = await finalshell.generateLicense(values.machineCode);
      setAuthorizationCodes(data);
    } catch (error) {
      console.error('Generate license failed:', error);
    } finally {
      setLoading(false);
    }
  };

  // Copy to clipboard
  const copyToClipboard = (key: string, text: string) => {
    copyAndManageState(
      key,
      text,
      copying,
      setCopying,
      notification,
      t('common.copied'),
      t('common.copyFail')
    );
  };

  const breadcrumbs = [
    {
      path: '/',
      breadcrumbName: t('nav.home'),
    },
    {
      path: '',
      breadcrumbName: t('nav.finalshell'),
    },
  ];

  return (
    <div>
      <PageHeader
        title={t('finalshell.title')}
        subTitle={t('finalshell.subTitle')}
        breadcrumbs={breadcrumbs}
      />

      <Paragraph>
        {t('finalshell.description')}
      </Paragraph>

      <UsageNotice
        message={t('finalshell.usageNotice')}
        description={t('finalshell.warningDescription')}
      />

      <FormWrapper>
        <Form form={form} onFinish={handleGenerateLicense} layout="vertical">
          <Row gutter={16}>
            <Col xs={24} sm={12}>
              <Form.Item
                name="machineCode"
                label={t('finalshell.machineCode')}
                rules={[{ required: true, message: t('finalshell.machineCodeRequired') }]}
              >
                <Input placeholder={t('finalshell.enterMachineCode')} />
              </Form.Item>
            </Col>
          </Row>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              {t('finalshell.generateButton')}
            </Button>
          </Form.Item>
        </Form>
      </FormWrapper>

      {authorizationCodes.length > 0 && (
        <FormCard title={t('finalshell.registrationSuccess')}>
          <Row gutter={16}>
            {authorizationCodes.map((code, index) => {
              const versionLabelKey = versionLabels[index] || '';
              const copyKey = `code-${index}`;

              return (
                <Col key={index} xs={24} sm={12} md={12} style={{marginBottom: 16}}>
                  <CodeLabel>{t(versionLabelKey)}</CodeLabel>
                  <CodeContainer>
                    <AuthorizationCodeContainer>
                      {code}
                      <ButtonContainer>
                        <CopyButton
                          size="small"
                          type="primary"
                          ghost
                          icon={copying[copyKey] ? <CheckOutlined /> : <CopyOutlined />}
                          onClick={() => copyToClipboard(copyKey, code)}
                        />
                      </ButtonContainer>
                    </AuthorizationCodeContainer>
                  </CodeContainer>
                </Col>
              );
            })}
          </Row>
        </FormCard>
      )}

      <FormCard title={t('finalshell.instructionsTitle')}>
        <StepItem>
          <StepNumber>1</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step1')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>2</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step2')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>3</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step3')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>4</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step4')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>5</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step5')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>6</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step6')}</StepContent>
        </StepItem>
        <StepItem>
          <StepNumber>7</StepNumber>
          <StepContent>{t('finalshell.usageSteps.step7')}</StepContent>
        </StepItem>
      </FormCard>
    </div>
  );
};

export default FinalShell; 