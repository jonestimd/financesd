package schema

import "github.com/jonestimd/financesd/internal/model"

var getAllCompanies = model.GetAllCompanies
var getCompanyByID = model.GetCompanyByID
var getCompanyByName = model.GetCompanyByName
var addCompanies = model.AddCompanies

var getAllAccounts = model.GetAllAccounts
var getAccountByID = model.GetAccountByID
var getAccountsByName = model.GetAccountsByName

var getAllCategories = model.GetAllCategories

var getAllGroups = model.GetAllGroups

var getAllPayees = model.GetAllPayees

var getAllSecurities = model.GetAllSecurities
var getSecurityByID = model.GetSecurityByID
var getSecurityBySymbol = model.GetSecurityBySymbol

var getAccountTransactions = model.GetTransactions
