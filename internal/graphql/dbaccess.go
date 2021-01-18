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
