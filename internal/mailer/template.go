package mailer

const forgotPasswordTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial; background: #f0f0f0; padding: 20px; }
        .card { max-width: 500px; margin: auto; background: white; padding: 25px; border-radius: 10px; }
        .code-box {
            font-size: 24px;
            font-weight: bold;
            background: #222;
            color: #fff;
            padding: 15px;
            text-align: center;
            border-radius: 6px;
            letter-spacing: 4px;
        }
    </style>
</head>
<body>
    <div class="card">
        <h2>Password Reset Request</h2>
        <p>Hello {{.Name}},</p>
        <p>Use the code below to reset your password:</p>

        <div class="code-box">{{.Code}}</div>

        <p>This code expires in 15 minutes.</p>
    </div>
</body>
</html>
`

const welcomeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial; background: #f7f7f7; padding: 20px; }
        .card { max-width: 500px; background: #fff; padding: 25px; margin: auto; border-radius: 10px; }
        .title { font-size: 22px; font-weight: bold; color: #333; }
        .msg { font-size: 16px; color: #555; margin-top: 10px; }
    </style>
</head>
<body>
    <div class="card">
        <div class="title">Welcome, {{.Name}} 🎉</div>
        <p class="msg">
            We're excited to have you onboard.  
            Let us know if you need help getting started.
        </p>
    </div>
</body>
</html>
`

const notificationTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial; background: #f5f5f5; padding: 20px; }
        .card { max-width: 500px; margin: auto; background: #fff; padding: 25px; border-radius: 10px; }
        .title { font-size: 20px; font-weight: bold; }
        .msg { margin-top: 10px; color: #555; font-size: 15px; }
    </style>
</head>
<body>
    <div class="card">
        <div class="title">{{.Title}}</div>
        <div class="msg">{{.Message}}</div>
    </div>
</body>
</html>
`

const verifyEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8" />
    <title>Verify Your Email</title>
    <style>
        body {
            background: #f5f7fa;
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
        }

        .container {
            max-width: 480px;
            background: white;
            margin: 40px auto;
            border-radius: 10px;
            padding: 30px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.08);
        }

        h2 {
            color: #333;
            text-align: center;
            font-weight: 600;
        }

        p {
            color: #555;
            font-size: 15px;
            line-height: 1.6;
        }

        .button {
            display: block;
            width: fit-content;
            background: #4f46e5;
            color: white;
            padding: 12px 20px;
            border-radius: 6px;
            text-decoration: none;
            font-weight: bold;
            margin: 25px auto;
            text-align: center;
        }

        .footer {
            text-align: center;
            margin-top: 25px;
            color: #999;
            font-size: 12px;
        }

        .code-box {
            margin: 20px auto;
            padding: 12px;
            background: #f0f2ff;
            border-left: 4px solid #4f46e5;
            border-radius: 5px;
            width: fit-content;
            font-size: 14px;
            color: #333;
        }

        a {
            color: #4f46e5;
        }
    </style>
</head>
<body>

<div class="container">
    <h2>Email Verification</h2>

    <p>Hello <strong>{{.Name}}</strong>,</p>

    <p>
        Thank you for signing up! To complete your registration, please verify
        your email address by clicking the button below.
    </p>

    <a href="{{.Link}}" class="button">Verify Email</a>

    <p>If the button doesn’t work, copy and paste the link below into your browser:</p>

    <div class="code-box">
        {{.Link}}
    </div>

    <p>
        If you did not create this account, please disregard this message.
    </p>

    <div class="footer">
        © {{now}} CBT App — All rights reserved.
    </div>
</div>

</body>
</html>
`
