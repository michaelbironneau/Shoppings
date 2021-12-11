USE OpDb
GO
TRUNCATE TABLE Security.Principal 
GO
EXECUTE [Security].[uspAddPrincipal] 
   'TestUser'
  ,'TestPass'
GO
DELETE FROM App.[List]
GO
TRUNCATE TABLE App.StoreOrder 
GO
DELETE FROM App.Store 
GO
DELETE FROM App.Item 
GO 