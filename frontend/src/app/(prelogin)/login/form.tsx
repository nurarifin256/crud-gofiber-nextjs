"use client";

import {
  Box,
  Button,
  Container,
  FormControl,
  IconButton,
  InputAdornment,
  InputLabel,
  OutlinedInput,
  TextField,
} from "@mui/material";
import { IconEye, IconEyeOff } from "@tabler/icons-react";
import React, { useState } from "react";

const FormLogin = () => {
  const [showPassword, setShowPassword] = useState<boolean>(false);

  const handleShowPassword = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setShowPassword((prev) => !prev);
  };

  const handleMouseDownPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
  };

  const handleMouseUpPassword = (
    event: React.MouseEvent<HTMLButtonElement>
  ) => {
    event.preventDefault();
  };
  return (
    <Container>
      <form autoComplete="off">
        <Box sx={{ mt: 1 }} display="flex" flexDirection="column">
          {/* email */}
          <TextField
            id="email"
            label="Email"
            variant="outlined"
            type="email"
            size="small"
            name="email"
            sx={{ mt: 3 }}
          />

          {/* password */}
          <FormControl sx={{ mt: 3, mb: 2 }} variant="outlined">
            <InputLabel htmlFor="password" size="small">
              Password
            </InputLabel>
            <OutlinedInput
              id="password"
              label="Password"
              type={showPassword ? "text" : "password"}
              size="small"
              autoComplete="off"
              endAdornment={
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={(e) => handleShowPassword(e)}
                    onMouseDown={(e) => handleMouseDownPassword(e)}
                    onMouseUp={(e) => handleMouseUpPassword(e)}
                  >
                    {showPassword ? (
                      <IconEye stroke={2} />
                    ) : (
                      <IconEyeOff stroke={2} />
                    )}
                  </IconButton>
                </InputAdornment>
              }
            />
          </FormControl>

          {/* button */}
          <Button type="button" fullWidth variant="contained" sx={{ my: 2 }}>
            Login
          </Button>
        </Box>
      </form>
    </Container>
  );
};

export default FormLogin;
