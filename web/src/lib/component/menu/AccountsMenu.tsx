import React from 'react';
import {observer} from 'mobx-react-lite';
import {Link} from 'react-router-dom';
import {MenuItem, ListItemText} from '@material-ui/core';
import {CompanyModel} from '../../model/CompanyModel';
import ChildMenu from './ChildMenu';
import {AccountModel} from '../../model/AccountModel';

type IMenuAccountProps = Pick<AccountModel, 'id' | 'name' | 'company'>;

export const MenuAccount = observer<IMenuAccountProps>(({id, name, company}) => {
    const showCompany = company?.filteredAccounts.length === 1;
    return (
        <MenuItem component={Link} to={`/finances/account/${id}`}>
            <ListItemText secondary={showCompany ? name : null}>{showCompany ? company?.name : name}</ListItemText>
        </MenuItem>
    );
});

interface IProps {
    onBack: () => void;
    company: CompanyModel;
}
// TODO disable/highlight current account on transactions page
const AccountsMenu = observer<IProps>(({onBack, company}) => {
    return (
        <ChildMenu onBack={onBack} name={company.name}>
            {company.filteredAccounts.map((account) => <MenuAccount key={account.id} {...account} />)}
        </ChildMenu>);
});

export default AccountsMenu;
