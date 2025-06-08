"use client";

import { ConfirmProvider } from "material-ui-confirm";

export default function PreLoginLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>){
    return (
        <ConfirmProvider>
            
        </ConfirmProvider>
    )
}
