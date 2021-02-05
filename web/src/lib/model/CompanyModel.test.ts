import {newAccountModel, newCompanyModel} from "src/test/accountFactory";

describe('CompanyModel', () => {
    describe('filteredAccounts', () => {
        it('returns accounts where hide is false', () => {
            const company = newCompanyModel({}, newAccountModel(), newAccountModel());
            jest.spyOn(company.accounts[0], 'hide', 'get').mockReturnValue(true);
            jest.spyOn(company.accounts[1], 'hide', 'get').mockReturnValue(false);

            expect(company.filteredAccounts).toEqual(company.accounts.slice(1));
        });
    });
});
