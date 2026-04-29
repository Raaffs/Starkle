import { Box,CardMedia,IconButton,useTheme } from "@mui/material"
import { useContext } from "react"
import { ColorModeContext,tokens } from "../../themes"
import  InputBase from "@mui/material/InputBase"
import  LightModeOutlinedIcon  from "@mui/icons-material/LightModeOutlined"
import  DarkModeOutlinedIcon  from "@mui/icons-material/DarkModeOutlined"
import  SearchIcon  from "@mui/icons-material/Search"
import  SettingsOutlinedIcon from "@mui/icons-material/SettingsOutlined"
import PersonOutlinedIcon from "@mui/icons-material/PersonOutlined"
import Header from "../../components/Header"
import logo from '../../assets/images/Proof.png'
const Topbar=()=>{
    const theme=useTheme();
    const colors=tokens(theme.palette.mode);
    const colorMode=useContext(ColorModeContext)
    return(
    <Box display="flex" justifyContent="space-between" p={3} >
      <Header title="Starkle" subtitle="Immutable"></Header>
      <Box/>
      <Box alignItems="center">
        <IconButton onClick={colorMode.toggleColorMode}>
          {theme.palette.mode === "dark" ? (
            <DarkModeOutlinedIcon />
          ) : (
            <LightModeOutlinedIcon />
          )}
        </IconButton>
        <IconButton>
        </IconButton>
        <IconButton>
          <PersonOutlinedIcon />
        </IconButton>
      </Box>
    </Box>
    ) 
}
export default Topbar