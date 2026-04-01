'use client';

import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Download, Upload, FileSpreadsheet, CheckCircle, XCircle, AlertCircle } from 'lucide-react';

interface BulkUploadModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

interface ValidationResult {
  row: number;
  status: 'success' | 'error' | 'warning';
  message: string;
}

export function BulkUploadModal({ open, onOpenChange }: BulkUploadModalProps) {
  const [file, setFile] = useState<File | null>(null);
  const [uploading, setUploading] = useState(false);
  const [validationResults, setValidationResults] = useState<ValidationResult[]>([]);
  const [uploadComplete, setUploadComplete] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) {
      setFile(selectedFile);
      setValidationResults([]);
      setUploadComplete(false);
    }
  };

  const handleDownloadTemplate = () => {
    // In a real implementation, this would download an actual Excel template
    const templateData = `Organization Name,Partner Type,Trade License,TIN,Contact Email,Contact Phone,Bank Account,Bank Name,Policy Types,Cashless Enabled,Cashless Limit,Discount Enabled,Discount Percentage,Service Locations,Acquisition Rate,Renewal Rate
LabAid Hospital Dhaka,HOSPITAL,TL-DHA-2026-123456,123456789012,contact@labaid.com,+8801712345678,1234567890123,Dutch Bangla Bank,"HEALTH,LIFE",TRUE,500000,FALSE,,Dhaka,20,10
Square Pharmacy Chittagong,PHARMACY,TL-CHI-2026-234567,234567890123,info@square.com,+8801812345678,2345678901234,Sonali Bank,HEALTH,FALSE,,TRUE,15,"Chittagong,Dhaka",15,8`;
    
    const blob = new Blob([templateData], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'partner_upload_template.csv';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  };

  const handleUpload = async () => {
    if (!file) return;

    setUploading(true);
    
    // Simulate file processing and validation
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Mock validation results
    const mockResults: ValidationResult[] = [
      { row: 1, status: 'success', message: 'LabAid Hospital Dhaka - Successfully validated' },
      { row: 2, status: 'success', message: 'Square Pharmacy Chittagong - Successfully validated' },
      { row: 3, status: 'error', message: 'Invalid TIN format - must be 12 digits' },
      { row: 4, status: 'warning', message: 'Duplicate trade license found - will update existing partner' },
      { row: 5, status: 'success', message: 'United Hospital Sylhet - Successfully validated' },
    ];

    setValidationResults(mockResults);
    setUploading(false);
    setUploadComplete(true);
  };

  const handleConfirmImport = () => {
    // In a real implementation, this would actually create the partners
    console.log('Importing validated partners...');
    onOpenChange(false);
    // Reset state
    setTimeout(() => {
      setFile(null);
      setValidationResults([]);
      setUploadComplete(false);
    }, 300);
  };

  const handleClose = () => {
    onOpenChange(false);
    // Reset state
    setTimeout(() => {
      setFile(null);
      setValidationResults([]);
      setUploadComplete(false);
    }, 300);
  };

  const successCount = validationResults.filter(r => r.status === 'success').length;
  const errorCount = validationResults.filter(r => r.status === 'error').length;
  const warningCount = validationResults.filter(r => r.status === 'warning').length;

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-3xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Bulk Upload Partners</DialogTitle>
          <DialogDescription>
            Upload an Excel file to add multiple partners at once. Download the template to see the required format.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-6">
          {/* Template Download */}
          <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
            <div className="flex items-center gap-3">
              <FileSpreadsheet className="h-8 w-8 text-muted-foreground" />
              <div>
                <div className="font-medium">Excel Template</div>
                <div className="text-sm text-muted-foreground">
                  Download the template with required columns and example data
                </div>
              </div>
            </div>
            <Button variant="outline" onClick={handleDownloadTemplate} className="gap-2">
              <Download className="h-4 w-4" />
              Download Template
            </Button>
          </div>

          {/* File Upload */}
          <div className="space-y-2">
            <label className="text-sm font-medium">Upload File</label>
            <div className="flex items-center gap-2">
              <input
                type="file"
                accept=".xlsx,.xls,.csv"
                onChange={handleFileChange}
                className="flex-1 text-sm file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-semibold file:bg-primary file:text-primary-foreground hover:file:bg-primary/90"
              />
              {file && (
                <Button
                  onClick={handleUpload}
                  disabled={uploading}
                  className="gap-2"
                >
                  {uploading ? (
                    <>
                      <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                      Validating...
                    </>
                  ) : (
                    <>
                      <Upload className="h-4 w-4" />
                      Validate
                    </>
                  )}
                </Button>
              )}
            </div>
            {file && (
              <div className="text-sm text-muted-foreground">
                Selected: {file.name} ({(file.size / 1024).toFixed(2)} KB)
              </div>
            )}
          </div>

          {/* Validation Results */}
          {validationResults.length > 0 && (
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <h3 className="font-medium">Validation Results</h3>
                <div className="flex gap-4 text-sm">
                  <span className="flex items-center gap-1 text-green-600">
                    <CheckCircle className="h-4 w-4" />
                    {successCount} Success
                  </span>
                  {warningCount > 0 && (
                    <span className="flex items-center gap-1 text-yellow-600">
                      <AlertCircle className="h-4 w-4" />
                      {warningCount} Warning
                    </span>
                  )}
                  {errorCount > 0 && (
                    <span className="flex items-center gap-1 text-red-600">
                      <XCircle className="h-4 w-4" />
                      {errorCount} Error
                    </span>
                  )}
                </div>
              </div>

              <div className="max-h-64 overflow-y-auto space-y-2 border rounded-lg p-4">
                {validationResults.map((result, index) => (
                  <div
                    key={index}
                    className={`flex items-start gap-2 text-sm p-2 rounded ${
                      result.status === 'success'
                        ? 'bg-green-50'
                        : result.status === 'error'
                        ? 'bg-red-50'
                        : 'bg-yellow-50'
                    }`}
                  >
                    {result.status === 'success' && (
                      <CheckCircle className="h-4 w-4 text-green-600 mt-0.5" />
                    )}
                    {result.status === 'error' && (
                      <XCircle className="h-4 w-4 text-red-600 mt-0.5" />
                    )}
                    {result.status === 'warning' && (
                      <AlertCircle className="h-4 w-4 text-yellow-600 mt-0.5" />
                    )}
                    <div>
                      <span className="font-medium">Row {result.row}:</span> {result.message}
                    </div>
                  </div>
                ))}
              </div>

              {errorCount === 0 && (
                <div className="p-4 bg-green-50 border border-green-200 rounded-lg">
                  <div className="flex items-center gap-2 text-green-800">
                    <CheckCircle className="h-5 w-5" />
                    <span className="font-medium">
                      All rows validated successfully! Ready to import {successCount} partners.
                    </span>
                  </div>
                </div>
              )}

              {errorCount > 0 && (
                <div className="p-4 bg-red-50 border border-red-200 rounded-lg">
                  <div className="flex items-center gap-2 text-red-800">
                    <XCircle className="h-5 w-5" />
                    <span className="font-medium">
                      Please fix {errorCount} error(s) before importing. Download the error report and correct the issues.
                    </span>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={handleClose}>
            Cancel
          </Button>
          {uploadComplete && errorCount === 0 && (
            <Button onClick={handleConfirmImport} className="gap-2">
              <CheckCircle className="h-4 w-4" />
              Confirm Import ({successCount} partners)
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
