package graphql

import "github.com/jonestimd/financesd/internal/model"

var getAllCompanies = model.GetAllCompanies
var getCompanyByID = model.GetCompanyByID
var getCompanyByName = model.GetCompanyByName
var getCompaniesByIDs = model.GetCompaniesByIDs

var getAllAccounts = model.GetAllAccounts
var getAccountByID = model.GetAccountByID
var getAccountsByCompanyIDs = model.GetAccountsByCompanyIDs
var getAccountsByName = model.GetAccountsByName

var getAllCategories = model.GetAllCategories

var getAllPayees = model.GetAllPayees

var getAllSecurities = model.GetAllSecurities
var getSecurityByID = model.GetSecurityByID
var getSecurityBySymbol = model.GetSecurityBySymbol
