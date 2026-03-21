import { Theme } from "@mui/material/styles";
export const textFieldSx = {
  padding: "5px",

  // Target the label when focused
  "& .MuiInputLabel-root.Mui-focused": {
    color: "#E94057", // Deep Coral
  },

  // Target the focused underline bar
  "& .MuiInput-underline:after": {
    borderBottomColor: "#E94057", // Deep Coral
    borderBottomWidth: "2px",
  },

  // Target the underline bar when hovered, before focus
  "& .MuiInput-underline:hover:not(.Mui-disabled):before": {
    borderBottomColor: "#F27121", // Lighter Orange
  },
};

export const btnstyle = {
  background: "linear-gradient(90deg, #E94057 10%, #F27121 90%)",
  padding: "8px 24px",
  color: "white",
  fontSize: "1rem",
  fontWeight: 600,
  borderRadius: "100px",
  boxShadow: "0 4px 10px 0 rgba(233, 64, 87, 0.4)",
  transition: "transform 0.2s ease-in-out",
};

export const menuItemStyle = {
  backgroundColor: "white",
  color: "#333",
  fontWeight: 500,
  fontSize: "1rem",
  transition: "0.15s ease",

  "&:hover": {
    backgroundColor: "rgba(233, 64, 87, 0.15)", // soft red tint
    color: "#E94057",
  },

  "&.Mui-selected": {
    backgroundColor: "#E94057", // strong red for selected
    color: "white",
  },

  "&.Mui-selected:hover": {
    backgroundColor: "#D7374F", // slightly darker red on hover
  },
};

export const DataGridSx = {
  width: "dynamic",
  maxWidth: "170vh",

  backgroundColor: "#e1e2fe",
  borderRadius: "18px",
  border: "1px solid #c3c6fd",
  boxShadow: "0 16px 40px rgba(0,0,0,0.08)",
  overflow: "hidden",
  color: "#0f172a",

  "& .MuiDataGrid-columnHeaders, & .MuiDataGrid-footerContainer": {
    background: "linear-gradient(135deg, #fb7185 0%, #f97316 100%)",
    color: "#ffffff",
    fontWeight: 700,
    border: "none",
  },

  "& .MuiDataGrid-footerContainer": {
    fontWeight: 500,
  },

  "& .MuiDataGrid-columnSeparator": {
    display: "none",
  },

  "& .MuiDataGrid-columnHeaders": {
    borderTopLeftRadius: "18px",
    borderTopRightRadius: "18px",
  },

  "& .MuiDataGrid-row": {
    backgroundColor: "#f1b9b7",
    color: "#475569",
    fontWeight: 400,
    transition: "transform 0.15s ease, box-shadow 0.15s ease",
  },

  "& .MuiDataGrid-row:nth-of-type(even)": {
    backgroundColor: "#f8fafc",
  },

  "& .MuiDataGrid-row:hover": {
    backgroundColor: "#f1f5f9",
    transform: "translateY(-1px)",
    boxShadow: "0 6px 16px rgba(15,23,42,0.08)",
  },

  "& .MuiDataGrid-row.Mui-selected": {
    backgroundColor: "#ffe4e6 !important",
  },

  "& .MuiDataGrid-cell": {
    borderBottom: "1px solid #f1f5f9",
    paddingY: 1.1,
  },

  "& .MuiCheckbox-root": {
    color: "#fde68a",
  },
  "& .MuiCheckbox-root.Mui-checked": {
    color: "#ffffff",
  },

  "& .MuiDataGrid-cell:focus, & .MuiDataGrid-columnHeader:focus": {
    outline: "none",
  },
};

export const DataGridDarkSx = {
  width: "dynamic",
  maxWidth: "170vh",

  backgroundColor: "#151632", 
  borderRadius: "18px",
  border: "1px solid #1e293b",
  boxShadow: "0 18px 45px rgba(0,0,0,0.55)",
  overflow: "hidden",
  color: "#e5e7eb",

  "& .MuiDataGrid-columnHeaders, & .MuiDataGrid-footerContainer": {
    background: "linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%)",
    color: "#ffffff",
    fontWeight: 700,
    border: "none",
  },

  "& .MuiDataGrid-footerContainer": {
    fontWeight: 500,
  },

  "& .MuiDataGrid-columnSeparator": {
    display: "none",
  },

  "& .MuiDataGrid-columnHeaders": {
    borderTopLeftRadius: "18px",
    borderTopRightRadius: "18px",
  },

  "& .MuiDataGrid-row": {
    backgroundColor: "#151632",
    color: "#cbd5f5",
    fontWeight: 400,
    transition: "transform 0.15s ease, box-shadow 0.15s ease",
  },

  "& .MuiDataGrid-row:nth-of-type(even)": {
    backgroundColor: "#020617",
  },

  /* classy dark hover */
  "& .MuiDataGrid-row:hover": {
    backgroundColor: "#2a2d64",
    transform: "translateY(-1px)",
    boxShadow: "inset 0 0 0 1px rgba(251,113,133,0.35)",
  },

  "& .MuiDataGrid-row.Mui-selected": {
    backgroundColor: "rgba(251,113,133,0.15) !important",
  },

  "& .MuiDataGrid-cell": {
    borderBottom: "1px solid #1e293b",
    paddingY: 1.1,
  },

  "& .MuiCheckbox-root": {
    color: "#fcd34d",
  },
  "& .MuiCheckbox-root.Mui-checked": {
    color: "#ffffff",
  },

  /* no light flash focus */
  "& .MuiDataGrid-cell:focus, & .MuiDataGrid-columnHeader:focus": {
    outline: "none",
  },
};

export const statBoxStyles = (theme: Theme, colors: any) => ({
  gridColumn: "span 3",
  backgroundColor:
    theme.palette.mode === "dark" ? colors.primary[800] :colors.blueAccent[900],
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
  boxShadow:
    theme.palette.mode === "dark"
      ? "0px 4px 12px rgba(0,0,0,0)"
      : "0px 4px 12px rgba(0,0,0,0.1)",
  borderRadius: "12px", // added rounded corners
});
//   

export const flexHeaderBoxStyles = (theme: Theme, colors: any) => ({
  display: "flex",
  justifyContent: "space-between",
  alignItems: "center",
  borderBottom: `4px solid ${colors.primary[500]}`,
  color: theme.palette.mode === "dark" ? colors.grey[300] : colors.grey[500],
  padding: "15px",

});