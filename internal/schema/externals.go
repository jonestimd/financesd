package schema

import (
	"github.com/jonestimd/financesd/internal/database"
	"github.com/jonestimd/financesd/internal/domain"
)

var getAllCompanies = domain.GetAllCompanies
var getCompanyByID = domain.GetCompanyByID
var getCompanyByName = domain.GetCompanyByName
var addCompanies = domain.AddCompanies
var updateCompanies = domain.UpdateCompanies
var deleteCompanies = database.DeleteCompanies

var getAllAccounts = domain.GetAllAccounts
var getAccountByID = domain.GetAccountByID
var getAccountsByName = domain.GetAccountsByName

var getAllCategories = database.GetAllCategories

var getAllGroups = database.GetAllGroups

var getAllPayees = database.GetAllPayees

var getAllSecurities = database.GetAllSecurities
var getSecurityByID = database.GetSecurityByID
var getSecurityBySymbol = database.GetSecurityBySymbol
var getAccountTransactions = domain.GetTransactions
var getTransactionsByIDs = domain.GetTransactionsByIDs
var insertTransactions = domain.InsertTransactions
var updateTransactions = domain.UpdateTransactions
var deleteTransactions = domain.DeleteTransactions
