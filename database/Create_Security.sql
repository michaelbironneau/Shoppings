CREATE SCHEMA Security
GO
CREATE TABLE Security.Principal
(
    PrincipalId SMALLINT NOT NULL IDENTITY(1,1) PRIMARY KEY,
    Username VARCHAR(64) NOT NULL,
    Passhash VARCHAR(512) NOT NULL,
    Salt UNIQUEIDENTIFIER,
    Token VARCHAR(255) NULL DEFAULT(NEWID())
    --FIXME: GUIDS are designed for uniqueness, not entropy!
);
CREATE UNIQUE INDEX ix_SecurityPrincipal on Security.Principal(Username);
GO
CREATE OR ALTER FUNCTION [Security].[udf_GetToken] (
	@pUsername VARCHAR(64),
	@pPassword VARCHAR(128)
)
RETURNS VARCHAR(255) AS
BEGIN
    DECLARE @ret VARCHAR(255);
    SELECT @ret = P.Token
    FROM Security.Principal P
    WHERE P.Username = @pUsername AND
        P.Passhash = HASHBYTES('SHA2_512', @pPassword + CAST(P.Salt AS NVARCHAR(36)))
    ;
    RETURN @ret
END;
GO
CREATE OR ALTER PROCEDURE [Security].[uspAddPrincipal]
    (
    @pUsername VARCHAR(64),
    @pPassword VARCHAR(128)
)
AS
BEGIN
    DECLARE @salt UNIQUEIDENTIFIER = NEWID();
    INSERT INTO [Security].[Principal]
        (Username, Passhash, Salt)
    VALUES
        (@pUsername, HASHBYTES('SHA2_512', @pPassword + CAST(@salt AS NVARCHAR(36))), @salt);
END;