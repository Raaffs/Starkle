import { useState } from "react";
import {
  Modal,
  Box,
  Typography,
  TextField,
  Button,
  IconButton,
  CircularProgress,
} from "@mui/material";
import CloseIcon from "@mui/icons-material/Close";
import { GenerateZKPOffline } from "../../../wailsjs/go/main/App";
import { textFieldSx, btnstyle } from "../../styles/styles"

interface ZKPModalProps {
  open: boolean;
  onClose: () => void;
}

const modalBoxSx = {
  position: "absolute" as const,
  top: "50%",
  left: "50%",
  transform: "translate(-50%, -50%)",
  width: 480,
  backgroundColor: "white",
  borderRadius: "18px",
  border: "1px solid #f1b9b7",
  boxShadow: "0 16px 40px rgba(233, 64, 87, 0.12)",
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
  color: "#888",
  "&:hover": { color: "#E94057" },
};

const feedbackBaseSx = {
  mt: 1.5,
  fontSize: "0.85rem",
  fontWeight: 500,
  px: 1.5,
  py: 0.75,
  borderRadius: "8px",
};

export default function ZKPModal({ open, onClose }: ZKPModalProps) {
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

  return (
    <Modal open={open} onClose={handleClose}>
      <Box sx={modalBoxSx}>

        {/* Header */}
        <Box sx={headerSx}>
          <Typography
            variant="h6"
            sx={{
              fontWeight: 700,
              background: "linear-gradient(90deg, #E94057 10%, #F27121 90%)",
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
          placeholder="e.g. age > 18, balance > 0, verified = true"
          helperText="Separate multiple constraints with commas"
          disabled={loading}
          sx={textFieldSx}
        />

        {/* Feedback */}
        {message && (
          <Typography
            sx={{
              ...feedbackBaseSx,
              color: "#166534",
              backgroundColor: "#dcfce7",
            }}
          >
            {message}
          </Typography>
        )}
        {error && (
          <Typography
            sx={{
              ...feedbackBaseSx,
              color: "#991b1b",
              backgroundColor: "#fee2e2",
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
            style={btnstyle}
            sx={{
              "&.Mui-disabled": {
                opacity: 0.45,
                color: "white !important",
                background:
                  "linear-gradient(90deg, #E94057 10%, #F27121 90%) !important",
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