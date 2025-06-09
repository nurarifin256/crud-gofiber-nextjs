"use client";

import AlertsProvider from "@/components/AlertContext";
import { ConfirmProvider } from "material-ui-confirm";

export default function PreLoginLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>){
    return (
        <ConfirmProvider>
            <AlertsProvider>
                {children}
            </AlertsProvider>
        </ConfirmProvider>
    )
}
