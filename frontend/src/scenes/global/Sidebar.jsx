import "react-pro-sidebar/dist/css/styles.css";
import { useState, useEffect } from "react"; 
import { IsApprovedInstitute } from "../../../wailsjs/go/main/App";
import { ProSidebar, Menu, MenuItem } from "react-pro-sidebar";
import { Box, IconButton, Typography, useTheme } from "@mui/material";
import { Link } from "react-router-dom";
import { tokens } from "../../themes";
import HomeOutlinedIcon from "@mui/icons-material/HomeOutlined";
import MenuOutlinedIcon from "@mui/icons-material/MenuOutlined";
import UploadFileOutlinedIcon from "@mui/icons-material/UploadFileOutlined";
import CancelOutlinedIcon from "@mui/icons-material/CancelOutlined";
import AddTaskOutlinedIcon from "@mui/icons-material/AddTaskOutlined";
import TimerOutlinedIcon from "@mui/icons-material/TimerOutlined";
import LogoutOutlinedIcon from "@mui/icons-material/LogoutOutlined";
import LoginOutlinedIcon from "@mui/icons-material/LoginOutlined";
import PersonAddOutlinedIcon from "@mui/icons-material/PersonAddOutlined";
import FileOpenOutlinedIcon from "@mui/icons-material/FileOpenOutlined";

const Item = ({ title, to, icon, selected, setSelected }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  
  let color = colors.grey[100]; 
  if (title === "Approved") {
    color = theme.palette.mode === "dark" ? "#4ade80" : "#16a34a";
  } else if (title === "Rejected") {
    color = theme.palette.mode === "dark" ? "#f87171" : "#dc2626";
  } else if (title === "Pending") {
    color = theme.palette.mode === "dark" ? "#60a5fa" : "#2563eb";
  }

  return (
    <MenuItem
      active={selected === title}
      style={{
        color: color,
        margin: "2px 0", // Reduced vertical margin to prevent jumpiness
      }}
      onClick={() => setSelected(title)}
      icon={icon}
    >
      <Typography sx={{ fontWeight: selected === title ? 700 : 400 }}>
        {title}
      </Typography>
      <Link to={to} />
    </MenuItem>
  );
};

// ... (imports and Item component remain exactly the same)

const Sidebar = ({ authStatus }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [selected, setSelected] = useState("Dashboard");
  const [isApproved, setIsApproved] = useState(null);

  useEffect(() => {
    if (authStatus) {
      IsApprovedInstitute()
        .then((result) => {
          setIsApproved(result);
        })
        .catch((err) => {
          console.error("Error fetching institute approval status:", err);
          setIsApproved(false);
        });
    } else {
      setIsApproved(false);
    }
  }, [authStatus]);

  if (isApproved === null && authStatus) {
    return (
      <Box sx={{ p: 4, color: colors.grey[100] }}>
        <Typography variant="h6">Loading Menu...</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        position: "sticky",
        display: "flex",
        height: "calc(100vh - 40px)",
        top: "20px",
        bottom: "20px",
        left: "15px",
        marginRight: "15px",
        zIndex: 10000,
        "& .pro-sidebar": {
          width: isCollapsed ? "80px" : "280px",
          minWidth: isCollapsed ? "80px" : "280px",
        },
        "& .pro-sidebar-inner": {
          background: `${
            theme.palette.mode === "dark" 
              ? "linear-gradient(180deg, #111827 0%, #0f172a 100%)" 
              : colors.blueAccent[900]
          } !important`,
          borderRadius: "24px",
          border: theme.palette.mode === "dark" 
            ? "1px solid rgba(255, 255, 255, 0.05)" 
            : "1px solid rgba(0, 0, 0, 0.05)",
          boxShadow: "0 10px 30px rgba(0,0,0,0.15)",
          overflow: "hidden",
        },
        "& .pro-icon-wrapper": {
          backgroundColor: "transparent !important",
        },
        "& .pro-inner-item": {
          padding: "8px 20px !important", // Balanced padding for centering
          transition: "background-color 0.2s ease-in-out, color 0.2s ease-in-out !important", // Smooth, non-bouncy transition
          margin: "4px 10px", // Standard margin for all items
          borderRadius: "12px",
        },
        // Hover effect only for non-active items
        "& .pro-menu-item:not(.active) .pro-inner-item:hover": {
          color: "#ffffff !important",
          background: `${colors.blueAccent[700]} !important`,
        },
        // Active tab styling
        "& .pro-menu-item.active": {
          background: theme.palette.mode === "dark"
            ? "linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)" // Professional Blue for Dark Mode
            : "linear-gradient(135deg, #fb7185 0%, #f97316 100%)", // Your Orange/Rose for Light Mode
          borderRadius: "12px",
          margin: "4px 10px",
          color: "#ffffff !important",
          // Disable hover styles when active
          "& .pro-inner-item:hover": {
            background: "transparent !important",
            cursor: "default",
          }
        },
        "& .ps-menu-root": {
           padding: "10px 0"
        }
      }}
    >
      <ProSidebar collapsed={isCollapsed}>
        <Menu iconShape="circle">
          <MenuItem
            onClick={() => setIsCollapsed(!isCollapsed)}
            icon={isCollapsed ? <MenuOutlinedIcon /> : undefined}
            style={{
              margin: "15px 0 25px 0",
              color: colors.grey[100],
            }}
          >
            {!isCollapsed && (
              <Box
                display="flex"
                justifyContent="space-between"
                alignItems="center"
                ml="15px"
              >
                <Typography
                  variant="h3"
                  color={theme.palette.mode === "dark" ? colors.greenAccent[300] : colors.greenAccent[300]}
                  sx={{ 
                    fontFamily: "'Poppins', sans-serif", 
                    fontWeight: 800,
                    letterSpacing: "1px"
                  }}
                >
                  Starkle
                </Typography>
                <IconButton onClick={() => setIsCollapsed(!isCollapsed)}>
                  <MenuOutlinedIcon sx={{ color: colors.grey[100] }} />
                </IconButton>
              </Box>
            )}
          </MenuItem>

          <Box paddingLeft={isCollapsed ? undefined : "5px"}>
            {authStatus && (
              <>
                <Typography
                  variant="h6"
                  color={colors.grey[400]}
                  sx={{ m: "20px 0 8px 25px", fontSize: "0.75rem", letterSpacing: "1.2px", textTransform: "uppercase" }}
                >
                  Overview
                </Typography>
                <Item title="Dashboard" to="/dashboard" icon={<HomeOutlinedIcon />} selected={selected} setSelected={setSelected} />

                <Typography
                  variant="h6"
                  color={colors.grey[400]}
                  sx={{ m: "25px 0 8px 25px", fontSize: "0.75rem", letterSpacing: "1.2px", textTransform: "uppercase" }}
                >
                  Documents
                </Typography>
                <Item title="Approved" to="/documents/approved" icon={<AddTaskOutlinedIcon />} selected={selected} setSelected={setSelected} />
                <Item title="Rejected" to="/documents/rejected" icon={<CancelOutlinedIcon />} selected={selected} setSelected={setSelected} />
                <Item title="Pending" to="/documents/pending" icon={<TimerOutlinedIcon />} selected={selected} setSelected={setSelected} />

                <Typography
                  variant="h6"
                  color={colors.grey[400]}
                  sx={{ m: "25px 0 8px 25px", fontSize: "0.75rem", letterSpacing: "1.2px", textTransform: "uppercase" }}
                >
                  Verify & Issue
                </Typography>
                {!isApproved && <Item title="Upload" to="/documents/upload" icon={<UploadFileOutlinedIcon />} selected={selected} setSelected={setSelected} />}
                {isApproved && <Item title="Issue" to="/documents/issue" icon={<FileOpenOutlinedIcon />} selected={selected} setSelected={setSelected} />}
              </>
            )}

            <Typography
              variant="h6"
              color={colors.grey[400]}
              sx={{ m: "25px 0 8px 25px", fontSize: "0.75rem", letterSpacing: "1.2px", textTransform: "uppercase" }}
            >
              Account
            </Typography>
            {authStatus ? (
              <Item title="Logout" to="/logout" icon={<LogoutOutlinedIcon />} selected={selected} setSelected={setSelected} />
            ) : (
              <>
                <Item title="Login" to="/" icon={<LoginOutlinedIcon />} selected={selected} setSelected={setSelected} />
                <Item title="New Account" to="/register" icon={<PersonAddOutlinedIcon />} selected={selected} setSelected={setSelected} />
              </>
            )}
          </Box>
        </Menu>
      </ProSidebar>
    </Box>
  );
};

export default Sidebar;