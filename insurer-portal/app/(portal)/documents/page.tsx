import { Suspense } from "react";

import { DocumentLibrary } from "@/components/document-library";

export default function DocumentsPage() {
  return (
    <Suspense fallback={null}>
      <DocumentLibrary />
    </Suspense>
  );
}
