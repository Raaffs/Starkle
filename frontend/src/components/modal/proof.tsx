import { useState } from "react";
import {
  Modal,
  Box,
  Typography,
  TextField,
  Button,
  IconButton,
  CircularProgress,
  useTheme,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { GenerateZKPOffline } from "../../../wailsjs/go/main/App";
import { textFieldSx, btnstyle } from "../../styles/styles";

interface ZKPModalProps {
  open: boolean;
  onClose: () => void;
}

export default function ZKPModal({ open, onClose }: ZKPModalProps) {
  const theme = useTheme();
  const isDarkMode = theme.palette.mode === "dark";

  const [constraints, setConstraints] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleClose = () => {
    setConstraints("");
    setMessage("");
    setError("");
    onClose();
  };

  const handleGenerate = () => {
    const constraintArray = constraints
      .split(",")
      .map((c) => c.trim())
      .filter(Boolean);

    setLoading(true);
    setMessage("");
    setError("");

    Promise.resolve()
      .then(() => GenerateZKPOffline("automatic", constraintArray))
      .then(() => {
        setMessage("Proof generated successfully.");
        setLoading(false);
      })
      .catch((err: unknown) => {
        setError(String(err));
        setLoading(false);
      });
  };

  // Dynamic Styles
  const modalBoxSx = {
    position: "absolute" as const,
    top: "50%",
    left: "50%",
    transform: "translate(-50%, -50%)",
    width: 480,
    backgroundColor: isDarkMode ? "#0f172a" : "white",
    backgroundImage: isDarkMode 
      ? "linear-gradient(180deg, rgba(30, 58, 138, 0.1) 0%, rgba(15, 23, 42, 0) 100%)" 
      : "none",
    borderRadius: "24px",
    border: isDarkMode ? "1px solid rgba(255, 255, 255, 0.1)" : "1px solid #f1b9b7",
    boxShadow: isDarkMode 
      ? "0 20px 50px rgba(0, 0, 0, 0.5)" 
      : "0 16px 40px rgba(233, 64, 87, 0.12)",
    p: 4,
    outline: "none",
  };

  const headerSx = {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    mb: 3,
  };

  const closeButtonSx = {
    color: isDarkMode ? "#94a3b8" : "#888",
    "&:hover": { color: isDarkMode ? "#3b82f6" : "#E94057" },
  };

  const feedbackBaseSx = {
    mt: 1.5,
    fontSize: "0.85rem",
    fontWeight: 500,
    px: 1.5,
    py: 0.75,
    borderRadius: "8px",
  };

  return (
    <Modal open={open} onClose={handleClose}>
      <Box sx={modalBoxSx}>
        {/* Header */}
        <Box sx={headerSx}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 800,
              letterSpacing: "0.5px",
              background: isDarkMode 
                ? "linear-gradient(90deg, #60a5fa 10%, #3b82f6 90%)"
                : "linear-gradient(90deg, #E94057 10%, #F27121 90%)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
            }}
          >
            Generate ZKP
          </Typography>
          <IconButton onClick={handleClose} size="small" sx={closeButtonSx}>
            <CloseIcon fontSize="small" />
          </IconButton>
        </Box>

        {/* Input */}

<TextField
  fullWidth
  label="Proof Constraints"
  variant="standard"
  value={constraints}
  onChange={(e) => setConstraints(e.target.value)}
  placeholder="e.g. age > 18, balance > 0"
  helperText="Separate multiple constraints with commas"
  disabled={loading}
  sx={{
    ...textFieldSx,
    // Label color
    "& .MuiInputLabel-root": {
      color: isDarkMode ? "#94a3b8" : "inherit",
    },
    "& .MuiInputLabel-root.Mui-focused": {
      color: isDarkMode ? "#3b82f6" : "#E94057", // Blue in dark, Red/Orange in light
    },
    // Input text and bottom line
    "& .MuiInput-root": {
      color: isDarkMode ? "#f8fafc" : "inherit",
      "&:before": { 
        borderBottomColor: isDarkMode ? "rgba(255,255,255,0.2)" : "inherit" 
      },
      "&:hover:not(.Mui-disabled):before": {
        borderBottomColor: isDarkMode ? "#60a5fa" : "inherit",
      },
      "&:after": {
        borderBottomColor: isDarkMode ? "#3b82f6" : "#E94057", // The focused "growing" line
      },
    },
    // Helper text
    "& .MuiFormHelperText-root": {
      color: isDarkMode ? "#64748b" : "inherit",
    }
  }}
/>
        {/* Feedback */}
        {message && (
          <Typography
            sx={{
              ...feedbackBaseSx,
              color: isDarkMode ? "#4ade80" : "#166534",
              backgroundColor: isDarkMode ? "rgba(22, 101, 52, 0.2)" : "#dcfce7",
            }}
          >
            {message}
          </Typography>
        )}
        {error && (
          <Typography
            sx={{
              ...feedbackBaseSx,
              color: isDarkMode ? "#f87171" : "#991b1b",
              backgroundColor: isDarkMode ? "rgba(153, 27, 27, 0.2)" : "#fee2e2",
            }}
          >
            {error}
          </Typography>
        )}

        {/* Generate Button */}
        <Box sx={{ display: "flex", justifyContent: "flex-end", mt: 4 }}>
          <Button
            onClick={handleGenerate}
            disabled={loading || !constraints.trim()}
            style={isDarkMode ? {} : btnstyle}
            sx={{
              textTransform: "none",
              fontWeight: 700,
              borderRadius: "12px",
              padding: "8px 24px",
              color: "white",
              background: isDarkMode
                ? "linear-gradient(45deg, #1e3a8a 30%, #3b82f6 90%)"
                : "linear-gradient(90deg, #E94057 10%, #F27121 90%)",
              boxShadow: isDarkMode 
                ? "0 4px 14px 0 rgba(0, 118, 255, 0.39)" 
                : "0 4px 14px 0 rgba(233, 64, 87, 0.39)",
              "&:hover": {
                background: isDarkMode
                  ? "linear-gradient(45deg, #1e40af 30%, #2563eb 90%)"
                  : "linear-gradient(90deg, #d83a4f 10%, #e1691e 90%)",
              },
              "&.Mui-disabled": {
                opacity: 0.45,
                color: "white !important",
                background: isDarkMode
                  ? "linear-gradient(45deg, #1e3a8a 30%, #3b82f6 90%) !important"
                  : "linear-gradient(90deg, #E94057 10%, #F27121 90%) !important",
              },
            }}
          >
            {loading ? (
              <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
                <CircularProgress size={16} sx={{ color: "white" }} />
                Generating...
              </Box>
            ) : (
              "Generate"
            )}
          </Button>
        </Box>
      </Box>
    </Modal>
  );
}