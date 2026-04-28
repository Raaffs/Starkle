//warning: this is all vibe coded
//it works but I've zero clue how. 
//at least it's not anything TOO important that breaks data integrity
import React, { useState, useEffect, ChangeEvent } from "react";
import {
  Box,
  Card,
  CardContent,
  Typography,
  TextField,
  Button,
  Grid,
  IconButton,
} from "@mui/material";
import { useTheme } from "@mui/material/styles"; 
import { IsApprovedInstitute, IssueCertificate } from "../../../wailsjs/go/main/App";
import { 
  FileText, 
  User, 
  Calendar, 
  Hash, 
  Home, 
  Info, 
  PlusCircle, // Added for the Add Field button
  Trash2, // Added for removing a custom field
  Type, // Icon for custom text field
} from "lucide-react"; 

import { tokens } from '../../themes'; 

// 1. UPDATED: CertificateData now uses a Record<string, any> union type to allow custom keys
export interface CertificateData {
  certificateName: string;
  publicAddress:string;
  name: string;
  address: string;
  age:  string;
  birthDate: string;
  uniqueId: string;
  [key: string]: any; // Allows any additional custom properties
}

interface IssueCardProps {
  data: CertificateData | null;
  viewTitle: string;
  onIssue: (cert:CertificateData) => void;
}

interface FieldConfig {
  label: string;
  name: keyof CertificateData | string; // Name can be a string for custom fields
  type: string;
  icon?: React.ElementType;
  gridWidth: 12 | 6 | 4;
  isCustom?: boolean; // Flag to identify custom fields
}

const IssueCard: React.FC<IssueCardProps> = ({ data: incomingData, viewTitle, onIssue }) => {
  const theme = useTheme();
  const colors = tokens(theme.palette.mode); 

  const [isApproved, setIsApproved] = useState<boolean>(false);
  
  // 2. NEW STATE: To manage dynamically added custom fields
  const [customFields, setCustomFields] = useState<FieldConfig[]>([]);
  // NEW STATE: To manage the input for the new custom field's label
  const [newCustomLabel, setNewCustomLabel] = useState<string>("");

  const emptyFields: CertificateData = {
    certificateName: "", name: "", address: "", age: "", birthDate: "", uniqueId: "",publicAddress:""
  };

  // State for all form data, including standard and custom fields
  const [data, setData] = useState<CertificateData>(emptyFields);

  // Standard (Fixed) Fields Configuration
  const standardFields: FieldConfig[] = [
    { label: "Certificate Name", name: "certificateName", icon: FileText, type: "text", gridWidth: 6 },
    { label: "Unique ID", name: "uniqueId", icon: Hash, type: "text", gridWidth: 6 },
    { label: "Public Address", name: "publicAddress", icon: Hash, type: "text", gridWidth: 12 },
    { label: "Name", name: "name", icon: User, type: "text", gridWidth: 4 },
    { label: "Age", name: "age", icon: Info, type: "number", gridWidth: 4 },
    { label: "Birth Date", name: "birthDate", icon: Calendar, type: "date", gridWidth: 4 },
    { label: "Address", name: "address", icon: Home, type: "text", gridWidth: 12 },
  ];

  // Combine standard and custom fields for rendering
  const allFields: FieldConfig[] = [...standardFields, ...customFields];


  useEffect(() => {
    IsApprovedInstitute().then((res: boolean) => {
      setIsApproved(res);

      if (!res && incomingData) {
        // When viewing data, we need to extract any custom fields that exist
        // This is simplified: standard fields are fixed, anything else is custom.
        const customDataKeys = Object.keys(incomingData).filter(
            k => !standardFields.map(f => f.name).includes(k)
        );
        
        // Create FieldConfig objects for the custom data found in incomingData
        const initialCustomFields: FieldConfig[] = customDataKeys.map(key => ({
            label: key.charAt(0).toUpperCase() + key.slice(1).replace(/([A-Z])/g, ' $1').trim(), // Simple label from camelCase key
            name: key,
            type: "text", // Assuming custom fields are text for simplicity
            icon: Type,
            gridWidth: 6,
            isCustom: true,
        }));
        
        setCustomFields(initialCustomFields);
        setData(incomingData);

      } else {
        setData(emptyFields);
        setCustomFields([]);
      }
    });
  }, [incomingData]);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (!isApproved) return;
    const { name, value } = e.target;
    
    // Ensure 'age' is stored as a number when it's a numeric field
    const updatedValue = (name === "age" && !isNaN(Number(value))) ? Number(value) : value;
    
    setData((prev) => ({
      ...prev,
      [name]: updatedValue.toString(), 
    }));
  };

  // 3. NEW HANDLER: To add a new custom field
  const handleAddField = () => {
    if (!newCustomLabel.trim()) return; // Prevent adding empty label
    
    const label = newCustomLabel.trim();
    // Create a unique camelCase name from the label for the state key
    const nameKey = label.toLowerCase().replace(/[^a-zA-Z0-9]+(.)/g, (match, chr) => chr.toUpperCase());
    
    const newField: FieldConfig = {
      label: label,
      name: nameKey,
      type: "text",
      icon: Type,
      gridWidth: 6, // Default to 50% width for custom fields
      isCustom: true,
    };

    setCustomFields((prev) => [...prev, newField]);
    setData((prev) => ({ ...prev, [nameKey]: "" })); // Initialize data for the new field
    setNewCustomLabel(""); // Reset input field
  };
  
  // NEW HANDLER: To remove a custom field
  const handleRemoveField = (fieldName: string) => {
      // Remove field from customFields state
      setCustomFields(prev => prev.filter(f => f.name !== fieldName));
      
      // Remove field's data from data state
      setData(prev => {
          const newState = { ...prev };
          delete newState[fieldName];
          return newState;
      });
  }


  return (
    <Card
      elevation={8}
      sx={{
        maxWidth: 750, 
        width: "100%",
        borderRadius: 3,
        margin: "16px auto", 
        background: `${theme.palette.mode==="dark" ? 'black' : 'white'} !important`,
      }}
    >
      <CardContent>
        <Box
          sx={{
            display: "flex",
            alignItems: "center",
            pb: 2,
            mb: 2,
          
            borderBottom: `2px solid ${theme.palette.primary.main}`, 
          }}
        >
          <FileText
            size={32}
            color={theme.palette.secondary.main} 
            style={{ marginRight: "16px" }}
          />
          <Typography variant="h5" color="secondary" fontWeight="bold">
            {viewTitle}
          </Typography>
        </Box>

        {isApproved ? (
          // Form mode (editable)
          <Grid container spacing={3}> 
            {/* 4. RENDER ALL FIELDS (Standard and Custom) */}
            {allFields.map((f) => (
              <Grid item xs={12} sm={f.gridWidth} key={f.name}> 
                <Box sx={{ 
                  display: 'flex', 
                  alignItems: 'flex-start' 

                  }}>
                    <TextField
                      fullWidth
                      variant="outlined"
                      label={f.label}
                      // Use the field's name property for the key
                      name={f.name as string}
                      type={f.type}
                      // Access the data dynamically
                      value={data[f.name as string] || ""} 
                      onChange={handleChange}
                      InputLabelProps={{
                        shrink: true,
                        sx: { color: `${theme.palette.mode=="dark"?colors.greenAccent[500]:colors.greenAccent[400]}`  }, 
                      }}
                      InputProps={{
                        startAdornment: f.icon ? (
                          <f.icon
                            size={18}
                            color={theme.palette.primary.main} 
                            style={{ marginRight: "8px" }}
                          />
                        ) : undefined,
                      }}
                    />
                    {/* Add a remove button ONLY for custom fields */}
                    {f.isCustom && (
                        <IconButton 
                            color="error"
                            onClick={() => handleRemoveField(f.name as string)}
                            sx={{ mt: 1, ml: 1, p: '8px' }}
                            aria-label={`Remove ${f.label}`}
                        >
                            <Trash2 size={20} />
                        </IconButton>
                    )}
                </Box>
              </Grid>
            ))}

            {/* NEW SECTION: Add Custom Field UI */}
            <Grid item xs={12}>
                <Box sx={{ 
                    display: 'flex', 
                    alignItems: 'center', 
                    mt: 2, 
                    p: 2,
                    border: `1px dashed ${theme.palette.divider}`,
                    borderColor:`${theme.palette.mode=="dark"?colors.greenAccent[600]:colors.greenAccent[400]}`,
                    borderRadius: 2
                }}>
                    <TextField
                        label="New Custom Field Name"
                        variant="outlined"
                        fullWidth
                        size="medium"
                        value={newCustomLabel}
                        onChange={(e) => setNewCustomLabel(e.target.value)}
                        onKeyDown={(e) => e.key === 'Enter' && handleAddField()}
                        sx={{ mr: 2 }}
                    />
                    <Button
                        variant="contained"
                        color="primary"
                        onClick={handleAddField}
                        startIcon={<PlusCircle size={20} />}
                        sx={{ minWidth: '150px' }}
                    >
                        Add Field
                    </Button>
                </Box>
            </Grid>
            {/* END NEW SECTION */}
          </Grid>
        ) : (
          // View-only mode (renders both standard and custom fields)
          <Grid container spacing={2}>
            {allFields.map((f) => (
              <Grid item xs={12} sm={f.gridWidth} key={f.name}> 
                <Box sx={{ mb: 1 }}>
                  <Typography variant="subtitle2" color="primary">
                    {f.label}
                  </Typography>
                  <Typography variant="body1" color="text.primary">
                    {String(data[f.name as string]) || "—"}
                  </Typography>
                </Box>
              </Grid>
            ))}
          </Grid>
        )}

        <Button
          fullWidth
          variant="contained"
          color={isApproved ? "secondary" : "primary"} 
          size="large"
          sx={{ mt: 3 }}
          onClick={()=>{
            console.log("Issuing certificate with data:", data); // Debug log to check the data being issued
            onIssue(data);
          }}
        >
          {isApproved ? "Issue Certificate" : ""}
        </Button>
      </CardContent>
    </Card>
  );
};

export default IssueCard;