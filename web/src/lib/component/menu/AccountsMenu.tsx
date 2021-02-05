import React from 'react';
import {observer} from 'mobx-react-lite';
import {Link} from 'react-router-dom';
import {MenuItem, ListItemText} from '@material-ui/core';
import {CompanyModel} from '../../model/CompanyModel';
import ChildMenu from './ChildMenu';
import {AccountModel} from '../../model/AccountModel';

type IMenuAccountProps = Pick<AccountModel, 'id' | 'name' | 'company'>;

export const MenuAccount: React.FC<IMenuAccountProps> = ({id, name, company}) => {
    const showCompany = company?.accounts.length === 1;
    return (
        <MenuItem>
            <ListItemText secondary={showCompany ? name : null}>
                <Link to={`/finances/account/${id}`}>{showCompany ? company?.name : name}</Link>
            </ListItemText>
        </MenuItem>
    );
};

interface IProps {
    onBack: () => void;
    company: CompanyModel;
}
// TODO disable/highlight current account on transactions page
// TODO close menu after selecting an account
const AccountsMenu = observer<IProps>(({onBack, company}) => {
    return (
        <ChildMenu onBack={onBack} name={company.name}>
            {company.filteredAccounts.map((account) => <MenuAccount key={account.id} {...account} />)}
        </ChildMenu>);
});

export default AccountsMenu;
