package model

import "github.com/jonestimd/financesd/internal/database"

var getAllCompanies = database.GetAllCompanies
var getCompanyByID = database.GetCompanyByID
var getCompaniesByIDs = database.GetCompaniesByIDs
var getCompanyByName = database.GetCompanyByName
var addCompany = database.AddCompany
var updateCompany = database.UpdateCompany
var deleteCompanies = database.DeleteCompanies

var getAllAccounts = database.GetAllAccounts
var getAccountByID = database.GetAccountByID
var getAccountsByName = database.GetAccountsByName
var getAccountsByCompanyIDs = database.GetAccountsByCompanyIDs
