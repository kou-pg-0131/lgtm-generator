import React, { useCallback, useMemo, useState } from 'react';
import Button from '@mui/material/Button';
import CardActions from '@mui/material/CardActions';
import CardContent from '@mui/material/CardContent';
import FormControlLabel from '@mui/material/FormControlLabel';
import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import TextField from '@mui/material/TextField';
import { styled } from '@mui/material/styles';
import { ReportType } from '~/types/report';
import { useSendReport } from './ReportHooks';
import ModalCard from '~/components/utils/ModalCard';
import LoadableButton from '~/components/utils/LoadableButton';

const StyledImage = styled('img')({});

type ReportFormProps = {
  lgtmId: string;
  open: boolean;
  onClose: () => void;
  imgSrc: string;
};

const ReportForm: React.VFC<ReportFormProps> = React.memo(props => {
  const { lgtmId, imgSrc, open, onClose } = props;

  const [type, setType] = useState<ReportType | null>(null);
  const [text, setText] = useState<string>('');

  const { sendReport, loading } = useSendReport();

  const isValid: boolean = useMemo(() => {
    if (!Object.values(ReportType).includes(type)) {
      return false;
    }
    if (text.length > 1000) {
      return false;
    }
    return true;
  }, [text.length, type]);

  const handleClose = useCallback(() => {
    if (loading) return;
    onClose();
  }, [loading, onClose]);

  const handleChangeType = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setType(e.currentTarget.value as ReportType);
    },
    [],
  );

  const handleChangeText = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setText(e.currentTarget.value);
    },
    [],
  );

  const handleSendReport = useCallback(() => {
    sendReport(lgtmId, type, text).then(() => {
      setText('');
      setType(null);
      onClose();
    });
  }, [lgtmId, onClose, sendReport, text, type]);

  return (
    <ModalCard open={open} onClose={handleClose}>
      <CardContent
        sx={{
          alignItems: 'center',
          display: 'flex',
          flexDirection: 'column',
          pt: 0,
        }}
      >
        <StyledImage
          src={imgSrc}
          alt='LGTM'
          sx={{
            border: '1px solid',
            borderColor: 'divider',
            mb: 2,
            maxHeight: 200,
            maxWidth: '100%',
          }}
        />
        <RadioGroup value={type || ''} onChange={handleChangeType}>
          <FormControlLabel
            value={ReportType.illegal}
            control={<Radio value={ReportType.illegal} />}
            label='法律違反 ( 著作権侵害、プライバシー侵害、名誉毀損等 )'
            disabled={loading}
          />
          <FormControlLabel
            value={ReportType.inappropriate}
            control={<Radio value={ReportType.inappropriate} />}
            label='不適切なコンテンツ'
            disabled={loading}
          />
          <FormControlLabel
            value={ReportType.other}
            control={<Radio value={ReportType.other} />}
            label='その他'
            disabled={loading}
          />
        </RadioGroup>
        <TextField
          fullWidth
          multiline
          placeholder='補足 ( 任意 )'
          disabled={loading}
          inputProps={{ maxLength: 1000 }}
          rows={5}
          onChange={handleChangeText}
          value={text}
        />
      </CardContent>
      <CardActions>
        <Button
          fullWidth
          color='secondary'
          onClick={handleClose}
          disabled={loading}
        >
          キャンセル
        </Button>
        <LoadableButton
          fullWidth
          color='primary'
          disabled={!isValid}
          loading={loading}
          onClick={handleSendReport}
        >
          送信
        </LoadableButton>
      </CardActions>
    </ModalCard>
  );
});

ReportForm.displayName = 'ReportForm';

export default ReportForm;
