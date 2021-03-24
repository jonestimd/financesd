import React from 'react';
import {observer} from 'mobx-react-lite';
import {Link, useHistory} from 'react-router-dom';
import {MenuItem, ListItemText} from '@material-ui/core';
import {CompanyModel} from '../../model/account/CompanyModel';
import ChildMenu from './ChildMenu';
import {AccountModel} from '../../model/account/AccountModel';

type IMenuAccountProps = Pick<AccountModel, 'id' | 'name' | 'company'>;

export const MenuAccount = observer<IMenuAccountProps>(({id, name, company}) => {
    const history = useHistory();
    const match = /^\/finances\/account\/(\d+)$/.exec(history.location.pathname);
    const currentAccount = (match && `${id}` === match[1]) ?? false;
    const showCompany = company?.filteredAccounts.length === 1;
    return (
        <MenuItem component={Link} to={`/finances/account/${id}`} selected={currentAccount}>
            <ListItemText secondary={showCompany ? name : null}>{showCompany ? company?.name : name}</ListItemText>
        </MenuItem>
    );
});

interface IProps {
    onBack: () => void;
    company: CompanyModel;
}

const AccountsMenu = observer<IProps>(({onBack, company}) => {
    return (
        <ChildMenu onBack={onBack} name={company.name}>
            {company.filteredAccounts.map((account) => <MenuAccount key={account.id} {...account} />)}
        </ChildMenu>);
});

export default AccountsMenu;
