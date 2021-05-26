package domain

import (
	"github.com/graphql-go/graphql"
	"github.com/jonestimd/financesd/internal/database"
)

var getAllCompanies = database.GetAllCompanies
var getCompanyByID = database.GetCompanyByID
var getCompaniesByIDs = database.GetCompaniesByIDs
var getCompanyByName = database.GetCompanyByName
var addCompany = database.AddCompany
var updateCompany = database.UpdateCompany

var getAllAccounts = database.GetAllAccounts
var getAccountByID = database.GetAccountByID
var getAccountsByName = database.GetAccountsByName
var getAccountsByCompanyIDs = database.GetAccountsByCompanyIDs

var getTransactions = database.GetTransactions
var getTransactionsByIDs = database.GetTransactionsByIDs
var updateTransaction = database.UpdateTransaction
var deleteTransactions = database.DeleteTransactions

var getDetailsByTxIDs = database.GetDetailsByTxIDs
var getDetailsByAccountID = database.GetDetailsByAccountID
var getRelatedTransactions = database.GetRelatedTransactions
var getRelatedTransactionsByAccountID = database.GetRelatedTransactionsByAccountID
var getRelatedDetailsByTxIDs = database.GetRelatedDetailsByTxIDs
var getRelatedDetailsByAccountID = database.GetRelatedDetailsByAccountID
var insertDetail = database.InsertDetail
var updateDetail = database.UpdateDetail
var validateDetails = database.ValidateDetails
var addOrUpdateTransfer = database.AddOrUpdateTransfer
var setTransferAmount = database.SetTransferAmount
var deleteDetails = database.DeleteDetails
var deleteRelatedDetails = database.DeleteRelatedDetails
var deleteTransactionDetails = database.DeleteTransactionDetails
var deleteTransfer = database.DeleteTransfer

var defaultResolveFn = graphql.DefaultResolveFn
