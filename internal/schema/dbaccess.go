package schema

import (
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/model"
)

var getAllCompanies = model.GetAllCompanies
var getCompanyByID = model.GetCompanyByID
var getCompanyByName = model.GetCompanyByName
var addCompanies = model.AddCompanies
var updateCompanies = model.UpdateCompanies
var deleteCompanies = database.DeleteCompanies

var getAllAccounts = model.GetAllAccounts
var getAccountByID = model.GetAccountByID
var getAccountsByName = model.GetAccountsByName

var getAllCategories = database.GetAllCategories

var getAllGroups = database.GetAllGroups

var getAllPayees = database.GetAllPayees

var getAllSecurities = database.GetAllSecurities
var getSecurityByID = database.GetSecurityByID
var getSecurityBySymbol = database.GetSecurityBySymbol
var getAccountTransactions = model.GetTransactions
var getTransactionsByIDs = model.GetTransactionsByIDs
var updateTransactions = model.UpdateTransactions
