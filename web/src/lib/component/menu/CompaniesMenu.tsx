import React from 'react';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../../store/RootStore';
import {Checkbox, Icon, ListItem, ListItemText, MenuItem} from '@material-ui/core';
import ChildMenu from './ChildMenu';
import {CompanyModel} from '../../model/account/CompanyModel';
import {MenuAccount} from './AccountsMenu';
import settingsStore from '../../store/settingsStore';
import {translate} from '../../i18n/localize';

interface IMenuCompanyProps {
    company: CompanyModel;
    selectCompany: (company: CompanyModel) => void
}

const MenuCompany = observer<IMenuCompanyProps>(({company, selectCompany}) => {
    if (company.accounts.length === 1) return <MenuAccount {...company.accounts[0]} />;
    return <MenuItem onClick={() => selectCompany(company)}><ListItemText>{company.name}</ListItemText><Icon>chevron_right</Icon></MenuItem>;
});

interface IProps {
    onBack: () => void;
    selectCompany: IMenuCompanyProps['selectCompany'];
}

const CompaniesMenu = observer<IProps>(({onBack, selectCompany}) => {
    const {accountStore} = React.useContext(RootStoreContext);
    const {hideClosedAccounts} = settingsStore;
    return (
        <ChildMenu onBack={onBack}>
            <ListItem divider>
                <ListItemText>{translate('menu.hideClosedAccounts')}</ListItemText>
                <Checkbox checked={hideClosedAccounts} onChange={() => settingsStore.hideClosedAccounts = !hideClosedAccounts} />
            </ListItem>
            {accountStore.accountsWithoutCompany.map((account) => <MenuAccount key={account.id} {...account}/>)}
            {accountStore.filteredCompanies.map((company) => <MenuCompany key={company.id} company={company} selectCompany={selectCompany} />)}
        </ChildMenu>);
});

export default CompaniesMenu;
