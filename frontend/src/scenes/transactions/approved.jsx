import { Box, Typography, useTheme, Button, Modal } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import { tokens } from "../../themes";
import Header from "../../components/Header";
import { Download, GetAcceptedDocs } from "../../../wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { DataGridSx, DataGridDarkSx } from "../../styles/styles";
import RemoveRedEyeSharpIcon from "@mui/icons-material/RemoveRedEyeSharp";
import DownloadForOfflineSharpIcon from "@mui/icons-material/DownloadForOfflineSharp";
import AutoAwesomeIcon from "@mui/icons-material/AutoAwesome";
import { ViewDigitalCertificate } from "../../../wailsjs/go/main/App";
import IssueCard from "../../components/cards/certificate";
import PopUp from "../../components/PopUp";
import ZKPModal from "../../components/modal/proof";

const ApprovedDocuments = () => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode);
  const [docs, setDocs] = useState([]);
  const [error, setError] = useState(null);
  const [modal, setModal] = useState(false);
  const [certificate, setCertificate] = useState(null);
  const [message, setMessage] = useState("");
  const [zkpModalOpen, setZkpModalOpen] = useState(false);
  const openZKPModal = () => setZkpModalOpen(true);
  const closeZKPModal = () => setZkpModalOpen(false);

  useEffect(() => {
    const getDocuments = () => {
      GetAcceptedDocs()
        .then((result) => {
          if (!result || result.length === 0) {
            setDocs([
              {
                ID: "",
                Requester: "",
                Verifier: "",
                ShaHash: "",
              },
            ]);
            setError("No Verified Documents");
          } else {
            const updatedDocs = result.map((doc) => {
              if (doc.IpfsAddress === "") {
                doc.IpfsAddress = "private";
              }
              return {
                ...doc,
                ShaHash: doc.ShaHash,
              };
            });
            setDocs(updatedDocs);
          }
        })
        .catch((err) => {
          setError(err.message);
        });
    };
    getDocuments();
  }, []);

  const handleView = (hash, institute, requester) => {
    setModal(true);
    ViewDigitalCertificate(hash, institute, requester)
      .then((data) => {
        console.log(data);
        setCertificate({
          certificateName: data.salted_fields.CertificateName.value,
          publicAddress: data.salted_fields.PublicAddress.value,
          name: data.salted_fields.Name.value,
          address: data.salted_fields.Address.value,
          age: data.salted_fields.Age.value,
          birthDate: data.salted_fields.BirthDate.value,
          uniqueId: data.salted_fields.UniqueID.value,
        });
        console.log("data: ", data);
      })
      .catch((err) => setError(err.message));
  };

  const columns = [
    { field: "Requester", headerName: "Requester", flex: 1 },
    { field: "Verifier", headerName: "Verifier", flex: 1 },
    { field: "ShaHash", headerName: "Hash", flex: 1 },
    {
      field: "view",
      headerName: "View",
      flex: 0.5,
      justifyContent: "center",
      renderCell: (params) => {
        return (
          <Box>
            <Button
              color="info"
              onClick={() => {
                handleView(
                  params.row.ShaHash,
                  params.row.Institute,
                  params.row.Requester,
                );
              }}
              sx={{
                minWidth: "auto",
                padding: "6px 8px",
                margin: 4,
                "&:hover": {
                  backgroundColor: "rgba(0, 0, 0, 0.1)",
                },
              }}
            >
              <RemoveRedEyeSharpIcon
                sx={{
                  color: `#2196F3`,
                }}
              />
            </Button>
          </Box>
        );
      },
    },
    {
      field: "download",
      headerName: "Download",
      flex: 0.5,
      justifyContent: "center",
      renderCell: (params) => {
        return (
          <Box>
            <Button
              color="info"
              onClick={() => {
                Promise.resolve()
                  .then(() => {
                    return Download(
                      params.row.ShaHash,
                      params.row.Institute,
                      params.row.Requester,
                    );
                  })
                  .then(() => {
                    setMessage("Downloaded successfully");
                  })
                  .catch((err) => {
                    setError(err);
                  });
              }}
              sx={{
                minWidth: "auto",
                padding: "6px 8px",
                margin: 4,
                "&:hover": {
                  backgroundColor: "rgba(0, 0, 0, 0.1)",
                },
              }}
            >
              <DownloadForOfflineSharpIcon
                sx={{
                  color: `#2196F3`,
                }}
              />
            </Button>
          </Box>
        );
      },
    },
  ];

  return (
    <Box
      m="20px"
      sx={{ width: "dynamic", maxWidth: "95%", justifyContent: "center" }}
    >
      {/* Top Header Section with Button Positioning */}
      <Box
        display="flex"
        justifyContent="space-between"
        alignItems="center"
        mb="20px"
      >
        <Header title="Approved Documents" />
        <Button
          onClick={openZKPModal}
          variant="contained"
          startIcon={<AutoAwesomeIcon />}
          sx={{
            // Logic for Light/Dark mode styling
            background:
              theme.palette.mode === "dark"
                ? "linear-gradient(45deg, #1e3a8a 30%, #3b82f6 90%)" // Your existing Dark Mode
                : "linear-gradient(135deg, #fb7185 0%, #f97316 100%)", // Matches your Light Mode Header

            color: "#ffffff",
            fontWeight: 700,
            padding: "10px 22px",
            borderRadius: "12px",
            textTransform: "none",
            fontSize: "0.95rem",
            boxShadow:
              theme.palette.mode === "dark"
                ? "0 4px 15px rgba(0,0,0,0.4)"
                : "0 4px 12px rgba(249, 115, 22, 0.25)",
            border:
              theme.palette.mode === "dark" ? "none" : "1px solid #c3c6fd",
            transition: "all 0.2s ease-in-out",

            "&:hover": {
              transform: "translateY(-2px)",
              filter: "brightness(1.05)",
              boxShadow:
                theme.palette.mode === "dark"
                  ? "0 6px 20px rgba(0,0,0,0.5)"
                  : "0 6px 16px rgba(251, 113, 133, 0.3)",
            },

            "&.Mui-disabled": {
              opacity: 0.5,
              color: "white !important",
            },
          }}
        >
          Auto Generate Proof
        </Button>{" "}
      </Box>

      {error && (
        <PopUp
          Error={error}
          Message=""
          onClose={() => {
            setError(null);
          }}
        />
      )}
      {message && (
        <PopUp
          Message={message}
          Error={null}
          onClose={() => {
            setMessage("");
          }}
        />
      )}

      {error && (
        <Typography
          color="error"
          align="center"
          style={{ marginBottom: "16px" }}
        >
          {error}
        </Typography>
      )}

      <Box
        m="10px 0 0 0"
        height="70vh"
        justifyContent="center"
        sx={{
          "& .MuiDataGrid-root": {
            border: "none",
          },
          "& .MuiDataGrid-cell": {
            borderBottom: "none",
            fontSize: "1.1rem",
          },
          "& .name-column--cell": {
            color: colors.greenAccent[300],
          },
          "& .MuiDataGrid-columnHeaders": {
            backgroundColor: colors.blueAccent[700],
            borderBottom: "none",
            fontSize: "1.2rem",
          },
          "& .MuiDataGrid-footerContainer": {
            borderTop: "none",
            backgroundColor: colors.blueAccent[900],
          },
          "& .MuiCheckbox-root": {
            color: `${colors.greenAccent[200]} !important`,
          },
        }}
      >
        <DataGrid
          columns={columns}
          rows={docs}
          getRowId={(row) => row.ID}
          sx={theme.palette.mode === "dark" ? DataGridDarkSx : DataGridSx}
        />
      </Box>

      <ZKPModal open={zkpModalOpen} onClose={closeZKPModal} />

      <Modal
        onClose={() => {
          setModal(false);
        }}
        open={modal}
        sx={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          p: 2,
        }}
      >
        <Box
          sx={{
            backgroundColor: "white",
            borderRadius: "18px",
            overflow: "hidden",
            background: `${
              theme.palette.mode === "dark" ? "black" : "transparent"
            } !important`,
          }}
        >
          <IssueCard
            data={certificate}
            viewTitle="Digital Certificate"
            onIssue={() => {}}
          />
        </Box>
      </Modal>
    </Box>
  );
};

export default ApprovedDocuments;
