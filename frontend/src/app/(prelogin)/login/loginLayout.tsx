import { Box, Card, Grid, Typography } from '@mui/material';
import styles from "./page.module.css";
import Image from 'next/image';
import React from 'react'

const LoginLayout = ({
    children,
    title,
    subtitle,
}: Readonly<{
    children: React.ReactNode;
    title: string;
    subtitle: string;
}>) => {
  return (
    <Card
        sx={{
            maxWidth: {md: 700, xs: 350},
            maxHeight: {md: 500, xs: 350},
            mx:"auto",
            borderRadius: 2,
            mb: 2,
        }}
    >
        <Grid container spacing={0}>
            {/* start form input */}
            <Grid size={{md:6, sm: 6, xs:12}} marginTop={2}>
                <Box marginTop={2} textAlign="center">
                    <Typography variant="h6" component="h6" fontWeight="bold">
                        {title}
                    </Typography>
                    <Typography variant="caption" component="h6" color="lightslategray" fontSize="0.7rem">
                        {subtitle}
                    </Typography>
                    {children}
                </Box>
            </Grid>
            {/* end form input */}

            {/* start image */}
            <Grid
                size={{md:6, sm:6, xs:12}}
                sx={{
                    backgroundColor: "#164acf",
                    display: {xs: "none", sm: "block", md: "block"},
                }}
            >
                <Box padding={3} display="flex" justifyContent="center" alignItems="center">
                    <Image
                        alt="illustration"
                        src="/images/login.jpg"
                        className={styles.illustration}
                        width={1000}
                        height={1000}
                        priority
                    />

                </Box>
            </Grid>
            {/* end image */}
        </Grid>
    </Card>
  )
}

export default LoginLayout