import { ChangeEvent, useEffect, useState } from "react";
import {
  Grid,
  Box,
  TextField,
  Typography,
  Button,
  Avatar,
  Link,
  Card,
  CardActions,
  InputLabel,
  Select,
  MenuItem,
} from "@mui/material";
import LockOutlinedIcon from "@mui/icons-material/LockOutlined";
// @ts-ignore
import { Login, Register } from "../../../wailsjs/go/main/App";
import { useNavigate } from "react-router-dom";
import { useTheme } from "@emotion/react";
import bg from "../../assets/images/Untitled.png";
import { btnstyle, menuItemStyle, textFieldSx } from "../../styles/styles";
import { PlatformOverviewCard } from "../../components/template/template";
function RegisterUser({ setAuthStatus }) {
  if (typeof Register !== "function") {
    console.log("Register function is not available. Please check your Wails setup.");
      setError("System not ready. Please wait a moment or restart.");
      return;
    }
  const theme = useTheme();
  let registerAsVerifier = false;
  const navigate = useNavigate();
  const [input, setInput] = useState({
    privateKey: "",
    username: "",
    password: "",
  });
  const [error, setError] = useState(null);
  const handleClick = () => {
    Register(
      input.privateKey,
      input.username,
      input.password,
      registerAsVerifier
    )
      .then(() => {
        setAuthStatus(true);
        navigate("/dashboard");
      })
      .catch((err) => {
        setError(err);
        console.log(err);
      });
  };

  const handleChange = (event) => {
    setInput({
      ...input,
      [event.target.name]: event.target.value,
    });
  };
  const setUserType = (Event) => {
    Event.target.value === "verifier"
      ? (registerAsVerifier = true)
      : (registerAsVerifier = false);
    console.log(registerAsVerifier);
  };
  return (
    <Box
      display="flex"
      padding={2}
      justifyContent="center"
      alignItems="center"
      sx={{
        backgroundColor: "transparent",
        backgroundSize: "100% 100%",
        backgroundImage: `url(${bg})`,
        backgroundRepeat: "no-repeat",
        width: "85vw",
        height: "86vh",
        borderRadius: "20px",
        margin: "10px",
      }}
    >
      <Card
        sx={{
          width: { xs: "90%", sm: "70%", md: "50%" }, // Responsive width
          minHeight: "450px",
          borderRadius: "16px",

          backgroundColor:
            theme.palette.mode == "light"
              ? "rgba(255, 255, 255, 0.9)"
              : "transparent",

          boxShadow: "0 8px 30px rgba(0, 0, 0, 0.2)",
          padding: "20px",
        }}
      >
        <Grid align="center" sx={{ mb: 2 }}>
          <Avatar
            sx={{
              m: 1,
              bgcolor: "#FF6F61", // lighter, brighter red
              color: "white", // text/icon stays visible
              top: "0px",
              boxShadow: "0 4px 12px rgba(255, 111, 97, 0.4)", // subtle glow for depth
            }}
          >
            <LockOutlinedIcon />
          </Avatar>
          <Typography
            variant="h4"
            component="h1"
            sx={{
              fontWeight: 600,
              color: theme.palette.mode == "light" ? "black" : "white",
            }}
          >
            Sign Up
          </Typography>
        </Grid>
        <CardActions sx={{ backgroundColor: "transparent" }}>
          <Box display="flex" flexDirection="column" width="100%">
            {error && (
              <Typography
                color="error"
                align="center"
                style={{ marginBottom: "16px" }}
              >
                {error}
              </Typography>
            )}
            <TextField
              label="Private Key"
              placeholder="Enter private key"
              variant="standard"
              fullWidth
              required
              name="privateKey"
              value={input.privateKey}
              onChange={handleChange}
              sx={textFieldSx}
            />
            <TextField
              label="Username"
              placeholder="Enter username"
              variant="standard"
              fullWidth
              required
              name="username"
              value={input.username}
              onChange={handleChange}
              sx={textFieldSx}
            />
            <TextField
              label="Password"
              placeholder="Enter password"
              type="password"
              variant="standard"
              fullWidth
              required
              name="password"
              value={input.password}
              onChange={handleChange}
              sx={textFieldSx}
            />
            <InputLabel
              sx={{mt:1.5}}
            >
             Register As
            </InputLabel>
            <Select
              variant="standard"
              id="demo-simple-select-standard"
              onChange={setUserType}
            >
              <MenuItem sx={{...menuItemStyle}} value="user">User</MenuItem>
              <MenuItem sx={{...menuItemStyle}} value="verifier">Verifier</MenuItem>
            </Select>
            <Box display="flex" justifyContent="center" alignItems="center">
              <Button
                type="submit"
                size="medium"
                variant="contained"
                onClick={handleClick}
                sx={{...btnstyle, mt:4,mb:4}}
              >
                Create Account
              </Button>
            </Box>
            <Box align="center" sx={{ mt: 2 }}>
              <Typography variant="body1" sx={{ color: "#555" }}>
                Already have an account?
              </Typography>
              <Typography
                component="a"
                href="/login"
                sx={{
                  color: "#E94057", // Matching the button/field color
                  fontWeight: 600,
                  textDecoration: "none",
                  "&:hover": { textDecoration: "underline" },
                }}
              >
                Log in
              </Typography>
              {/* <Redirect to="/register" /> - Replace with the Typography link above */}
            </Box>
          </Box>
        </CardActions>
      </Card>
      <PlatformOverviewCard/>
    </Box>
  );
}

export default RegisterUser;
