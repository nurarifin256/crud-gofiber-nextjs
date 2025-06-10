import React from 'react'
import LoginLayout from './loginLayout'
import FormLogin from './form'

const LoginPage = () => {
  return (
    <LoginLayout title="Login" subTitle="Please log in to continue">
        <FormLogin />
    </LoginLayout>
  )
}

export default LoginPage