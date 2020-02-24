@echo off
echo This bat file will install the AbsorbLMS.GoAws.Service and then start it.
echo Intended for developers usage, and not during deployment.
echo The service simulates a limited version of Amazoin SNS and SQS services.
echo View more details on GoAws project in "https://github.com/p4tin/goaws".
echo You need to provide three arguments to this bat for it to work.

echo Installing service comes next...
pause
echo Running command to install AbsorbLMS.GoAws.Service...
sc create newservice binpath=C:\code\github\goaws\goawssvc.exe
start /wait sc create AbsorbLMS.GoAws.Service binPath= "%~dp0\AbsorbLMS.GoAws.Service.exe" start= auto
echo Install command finished.
echo Starting service comes next...
pause
echo Running command to start AbsorbLMS.GoAws.Service...
start /wait net start AbsorbLMS.GoAws.Service
echo Start command finished.
goto :EOF