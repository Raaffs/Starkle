import { Box, Typography, useTheme, Button, Modal } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import { tokens } from "../../themes";
import Header from "../../components/Header";
import { Download, GetAcceptedDocs } from "../../../wailsjs/go/main/App";
import { useEffect, useState } from "react";
import { DataGridSx, DataGridDarkSx } from "../../styles/styles";
import RemoveRedEyeSharpIcon from "@mui/icons-material/RemoveRedEyeSharp";
import DownloadForOfflineSharpIcon from "@mui/icons-material/DownloadForOfflineSharp";
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
        const doc = docs.find((doc) => doc.ID === params.id);
        return (
          <Box>
            <Button
              color="info"
              onClick={() => {
                handleView(
                  params.row.ShaHash,
                  params.row.Institute,
                  params.row.Requester
                );
              }}
              sx={{
                // Ensures the button is compact for an icon-only use case
                minWidth: "auto",
                padding: "6px 8px",
                margin: 4,
                // Subtle background for better visibility against white
                // backgroundColor: 'rgba(0, 0, 0, 0.03)',
                // Use hardcoded color for hover border
                "&:hover": {
                  backgroundColor: "rgba(0, 0, 0, 0.1)",
                  // borderColor: '#64B5F6', // A light blue border on hover
                },
              }}
            >
              <RemoveRedEyeSharpIcon
                sx={{
                  // Professional Blue Color (e.g., MUI's standard info.main: #2196F3)
                  color: `linear-gradient(45deg, #00C6FF 30%, #0072FF 90%)`,
                  // fontSize: 20, // Adjusted size for visual balance
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
        const doc = docs.find((doc) => doc.ID === params.id);
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
                      params.row.Requester
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
                // Ensures the button is compact for an icon-only use case
                minWidth: "auto",
                padding: "6px 8px",
                margin: 4,
                // Subtle background for better visibility against white
                // backgroundColor: 'rgba(0, 0, 0, 0.03)',
                // Use hardcoded color for hover border
                "&:hover": {
                  backgroundColor: "rgba(0, 0, 0, 0.1)",
                  // borderColor: '#64B5F6', // A light blue border on hover
                },
              }}
            >
              <DownloadForOfflineSharpIcon
                sx={{
                  // Professional Blue Color (e.g., MUI's standard info.main: #2196F3)
                  color: `linear-gradient(45deg, #00C6FF 30%, #0072FF 90%)`,
                  // fontSize: 20, // Adjusted size for visual balance
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
       <Button
        onClick={openZKPModal}
        sx={{
              "&.Mui-disabled": {
                opacity: 0.45,
                color: "white !important",
                background:
                  "linear-gradient(90deg, #E94057 10%, #F27121 90%) !important",
              },
            }}
      >Auto Generate Proof
      </Button>

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
            setError(null);
          }}
        />
      )}

      <Header title="Approved Documents" />
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
        m="40px 0 0 0"
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
      >        <DataGrid
          columns={columns}
          rows={docs}
          getRowId={(row) => row.ID} // Use `ID` as a unique identifier
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
            // display: "flex",
            // gap: 4,
            // maxHeight: "92vh",
            overflow: "hidden",
            background: `${
              theme.palette.mode === "dark" ? "black" : "transparent"
            } !important`,

            // p: 4,
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
