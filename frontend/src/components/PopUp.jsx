import * as React from "react";
import { Snackbar, Alert, Box, LinearProgress } from "@mui/material";

const PopUp = ({ Message, Error, onClose }) => {
  const [open, setOpen] = React.useState(true);
  const [progress, setProgress] = React.useState(100);
  const duration = 6000;

  const handleClose = (event, reason) => {
    if (reason === "clickaway") return;
    setOpen(false);
    if (onClose) onClose();
  };

  React.useEffect(() => {
    if (open) {
      const startTime = Date.now();
      const timer = setInterval(() => {
        const elapsedTime = Date.now() - startTime;
        const remaining = Math.max(0, 100 - (elapsedTime / duration) * 100);
        setProgress(remaining);
      }, 10);

      return () => clearInterval(timer);
    }
  }, [open]);

  return (
    <Snackbar
      open={open}
      autoHideDuration={duration}
      onClose={handleClose}
      anchorOrigin={{ vertical: "bottom", horizontal: "right" }}
      sx={{
        mb: 2,
        mr: 2,
        width: "100%",
        maxWidth: "420px", // keeps it elegant, not stretched too much
      }}
    >
      <Box sx={{ width: "100%" }}>
        <Alert
          onClose={handleClose}
          severity={Error ? "error" : "success"}
          variant="filled"
          sx={{
            width: "100%",
            display: "flex",
            alignItems: "center",
            fontSize: "0.95rem",
            fontWeight: 500,
            borderRadius: "14px",
            boxShadow: "0px 6px 18px rgba(0,0,0,0.08)",
            backdropFilter: "blur(6px)",
            backgroundColor: Error
              ? "rgba(244, 67, 54, 0.92)"
              : "rgba(76, 175, 80, 0.92)",
            color: "#fff",
            position: "relative",
            overflow: "hidden",
            paddingBottom: "6px",

            "& .MuiAlert-icon": {
              fontSize: "22px",
              opacity: 0.9,
            },

            "& .MuiAlert-message": {
              width: "100%",
            },
          }}
        >
          <Box sx={{ width: "100%" }}>
            <Box sx={{ lineHeight: 1.4 }}>
              {Error ? Error : Message}
            </Box>

            <LinearProgress
              variant="determinate"
              value={progress}
              sx={{
                position: "absolute",
                bottom: 0,
                left: 0,
                width: "100%",
                height: "3px",
                backgroundColor: "rgba(255,255,255,0.25)",
                "& .MuiLinearProgress-bar": {
                  backgroundColor: "rgba(255,255,255,0.8)",
                  transition: "none",
                },
              }}
            />
          </Box>
        </Alert>
      </Box>
    </Snackbar>
  );
};

export default PopUp;