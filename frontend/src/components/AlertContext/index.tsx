import { subscribe } from "@/helpers/events";
import { EVENT_ERROR_API_GENERAL } from "@/helpers/restApi";
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContentText,
  Typography,
} from "@mui/material";
import {
  IconAlertCircleFilled,
  IconCircleCheckFilled,
  IconCircleXFilled,
} from "@tabler/icons-react";
import {
  createContext,
  JSX,
  ReactNode,
  useEffect,
  useMemo,
  useState,
} from "react";

type AlertContextType = {
  show: (
    title: string,
    type: string,
    content?: ReactNode | string,
    buttonLabel?: string,
    buttonLabelCancel?: string,
    buttonOnPress?: () => void
  ) => void;
};

type IconMapperType = {
  [key: string]: JSX.Element | null;
  success: JSX.Element;
  warning: JSX.Element;
  confirm: JSX.Element;
  danger: JSX.Element;
  info: JSX.Element;
};

const AlertContext = createContext<AlertContextType>({
  show: () => {},
});

const AlertsProvider = ({ children }: { children: ReactNode }) => {
  const IconMapper: IconMapperType = {
    success: <IconCircleCheckFilled width={130} height={130} color="#2BA04C" />,
    warning: <IconAlertCircleFilled width={130} height={130} color="#FFD503" />,
    confirm: <IconAlertCircleFilled width={130} height={130} color="#FFD503" />,
    danger: <IconCircleXFilled width={130} height={130} color="#B82222" />,
    info: <IconAlertCircleFilled width={130} height={130} color="#608DFF" />,
  };

  const [open, setOpen] = useState<boolean>(false);
  const [title, setTitle] = useState<string>("");
  const [buttonLabel, setButtonLabel] = useState<string>("Ok");
  const [buttonLabelCancel, setButtonLabelCancel] = useState<string>("Cancel");
  const [type, setType] = useState<string>("");
  const [isAuthorized, setIsAuthorized] = useState<boolean>(false);
  const [content, setContent] = useState<ReactNode | string>("");
  const [buttonOnPress, setButtonOnPress] = useState<() => void>(
    () => () => {}
  );

  const close = () => {
    setOpen(false);
    if (isAuthorized) location.href = "/login";
  };

  useEffect(() => {
    setIsAuthorized(false);
    subscribe(EVENT_ERROR_API_GENERAL, (data: CustomEvent) => {
      setTitle(data.detail.code);
      setContent(data.detail.message);
      setType("danger");
      setOpen(true);
    });
    subscribe(EVENT_ERROR_API_GENERAL, (data: CustomEvent) => {
      setTitle(data.detail.code);
      setContent(data.detail.message);
      setType("danger");
      setIsAuthorized(true);
      setOpen(true);
    });
  }, []);

  const show = useMemo(() => {
    return (
      title: string,
      type: string,
      content?: ReactNode | string,
      buttonLabel?: string,
      buttonLabelCancel?: string,
      buttonOnPress?: () => void
    ) => {
      setTitle(title);
      if (content) setContent(content);
      setType(type);
      if (buttonLabel) setButtonLabel(buttonLabel);
      if (buttonLabelCancel) setButtonLabelCancel(buttonLabelCancel);
      if (buttonOnPress) setButtonOnPress(() => buttonOnPress);
      setOpen(true);
    };
  }, [setTitle, setContent, setButtonLabel, setOpen]);

  const contextValue = useMemo(() => ({ show }), [show]);

  return (
    <AlertContext.Provider value={contextValue}>
      <Dialog
        fullWidth
        maxWidth="xs"
        open={open}
        onClose={close}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        sx={{ zIndex: 9999 }}
      >
        <Box
          display="flex"
          alignItems="center"
          justifyContent="center"
          marginTop="28px"
        >
          {IconMapper[type] || null}
        </Box>
        <Box
          display="flex"
          alignItems="center"
          justifyContent="center"
          marginTop="28px"
        >
          <Typography
            variant="h4"
            fontWeight={"medium"}
            color="primary"
            padding={0}
          >
            {title}
          </Typography>
        </Box>
        {content && (
          <Box
            display="flex"
            alignItems="center"
            justifyContent="center"
            marginTop="12px"
            marginBottom="28px"
          >
            <DialogContentText
              paddingX={2}
              id="alert-dialog-description"
              variant="body1"
              sx={{
                fontSize: 18,
                fontWeight: 400,
                whiteSpace: "pre-line",
                textAlign: "center",
              }}
            >
              {content}
            </DialogContentText>
          </Box>
        )}
        <DialogActions>
          <Box
            display="flex"
            justifyContent={type === "confirm" ? "end" : "center"}
            width="100%"
            sx={{
              padding: "24px",
            }}
            marginBottom={1}
          >
            {type === "confirm" && (
              <Button
                type="button"
                variant="outlined"
                color="secondary"
                onClick={close}
                sx={{
                  width: "138px",
                  height: "48px",
                  borderRadius: "15px",
                  marginRight: "16px",
                }}
              >
                {buttonLabelCancel}
              </Button>
            )}
            <Button
              type="button"
              variant="contained"
              sx={{
                width: "138px",
                height: "48px",
                borderRadius: "15px",
              }}
              onClick={
                type === "confirm"
                  ? () => {
                      buttonOnPress();
                      close();
                    }
                  : close
              }
            >
              {buttonLabel}
            </Button>
          </Box>
        </DialogActions>
      </Dialog>
      {children}
    </AlertContext.Provider>
  );
};
