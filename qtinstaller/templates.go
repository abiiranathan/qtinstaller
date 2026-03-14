package qtinstaller

const configTextTemplate = `
# Application display name
%s=

# Company name
%s=

# Description of what the software does.
%s=

# Release date
%s=2023-01-12

# Version number of the software
%s=0.1.0-1

# Package name for the software
%s=com.example.appname

# The path to the release Qt application
%s=%s

# Output name for the installer
%s=%s

# Name of the license
%s=MIT LICENCE

# Path to the license file
%s=license.txt

# Filename for a logo in PNG format used as QWizard::LogoPixmap.
%s=logo.png

# Filename for a custom installer icon (optional on Linux).
# On Windows use .ico format. On Linux use .png format.
%s=%s

`

const configXMLTemplate = `<?xml version="1.0" encoding="UTF-8"?> 
<Installer> 
    <Name>%s</Name> 
    <Version>%s</Version> 
    <Title>%s</Title> 
    <Publisher>%s</Publisher> 
    <StartMenuDir>%s</StartMenuDir> 
    <InstallerWindowIcon>%s</InstallerWindowIcon> 
    <InstallerApplicationIcon>%s</InstallerApplicationIcon> 
    <Logo>%s</Logo> 
    <TargetDir>@ApplicationsDir@/%s</TargetDir> 
    <ControlScript>controlscript.qs</ControlScript>
    <WizardDefaultWidth>700</WizardDefaultWidth> 
    <WizardDefaultHeight>500</WizardDefaultHeight> 
</Installer> 

`

const packageXMLTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<Package>
    <DisplayName>%s</DisplayName>
    <Description>%s</Description>
    <Version>%s</Version>
    <ReleaseDate>%s</ReleaseDate>
    <Licenses>
        <License name="%s" file="%s" />
    </Licenses>
    <Default>true</Default>
    <Script>installscript.qs</Script>
</Package>

`

const controlScriptTmpl = `
function Controller()
{
    installer.setDefaultPageVisible(QInstaller.TargetDirectory, true);
}
`

const installScriptTmpl = `
function Component(){}

Component.prototype.createOperations = function() 
{
    // call default implementation to actually install executable
    component.createOperations()
    if (systemInfo.productType === "windows") {
        component.addOperation("CreateShortcut", 
            "@TargetDir@/%s", 
            "@StartMenuDir@/%s",
            "workingDirectory=@TargetDir@", 
            "iconPath=@TargetDir@/%s",
            "iconId=0", 
            "description=Start Application")

        component.addOperation("CreateShortcut", 
            "@TargetDir@/%s", 
            "@DesktopDir@/%s",
            "workingDirectory=@TargetDir@", 
            "iconPath=@TargetDir@/%s",
            "iconId=0", 
            "description=Start Application")
    } else {
        component.addOperation("Execute", "chmod", "-R", "+x", "@TargetDir@",
            "UNDOEXECUTE", "true")
        component.addOperation("CreateDesktopEntry",
            "@HomeDir@/.local/share/applications/%s",
            "Type=Application\nExec=@TargetDir@/%s\nPath=@TargetDir@\nIcon=@TargetDir@/%s\nName=%s\nCategories=Utility;\nTerminal=false")
    } 
}

`
