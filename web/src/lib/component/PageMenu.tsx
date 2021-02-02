import React from 'react';
import {Link} from 'react-router-dom';
import MenuItem from '@material-ui/core/MenuItem';
import MenuList from '@material-ui/core/MenuList';
import {translate} from '../i18n/localize';
import {Icon, IconButton, ListItemText} from '@material-ui/core';
import {CompanyModel} from '../model/CompanyModel';
import {observer} from 'mobx-react-lite';
import {RootStoreContext} from '../store/RootStore';

const menuItems = [
    // ['menu.accounts', '/finances/'],
    ['menu.categories', '/finances/categories'],
    ['menu.groups', '/finances/groups'],
    ['menu.payees', '/finances/payees'],
    ['menu.securities', '/finances/securities'],
];

interface IProps {
    currentPage?: string;
}

const MenuAccount: React.FC<{id: string, name: React.ReactNode}> = ({id, name}) => (
    <MenuItem><ListItemText><Link to={`/finances/account/${id}`}>{name}</Link></ListItemText></MenuItem>
);

const MenuCompany: React.FC<{company: CompanyModel, callback: (c: CompanyModel) => void}> = ({company, callback}) => {
    if (company.accounts.length === 1) {
        const [{id, name}] = company.accounts;
        return <MenuAccount id={id} name={<span>{company.name}{' \u25BA '}{name}</span>} />;
    }
    return <MenuItem onClick={() => callback(company)}><ListItemText>{company.name}</ListItemText><Icon>chevron_right</Icon></MenuItem>;
};

const PageMenu = observer<IProps, HTMLDivElement>(({currentPage}) => {
    const {accountStore} = React.useContext(RootStoreContext);
    const [widths, setWidths] = React.useState<Record<string, number>>({});
    const [depth, setDepth] = React.useState(0);
    const [company, setCompany] = React.useState<CompanyModel | null>(null);
    const items = currentPage ? menuItems.filter(([key]) => currentPage !== key) : menuItems;
    const setWidth = (name: string, width: number) => {
        if (width > (widths[name] ?? 0)) setWidths({...widths, [name]: width});
    };
    const size = Object.values(widths).reduce((max, width) => width > max ? width : max, 0);
    const selectCompany = (company: CompanyModel) => {
        setCompany(company);
        setDepth(2);
    };
    return (
        <div className='sidebar-menu' style={{maxWidth: Math.max(size, 150)}}>
            <div className='menu' style={{width: Math.max(size*3, 450), translate: -depth*size}}>
                <MenuList ref={(el) => setWidth('top', el?.clientWidth ?? 0)}>
                    <MenuItem>
                        <ListItemText><Link to='/finances/'>Accounts</Link></ListItemText>
                        <IconButton size='small' onClick={() => setDepth(1)}><Icon>chevron_right</Icon></IconButton>
                    </MenuItem>
                    {items.map(([key, url]) => (<MenuItem key={key}><Link to={url}>{translate(key)}</Link></MenuItem>))}
                </MenuList>
                <MenuList ref={(el) => setWidth('companies', el?.clientWidth ?? 0)}>
                    <MenuItem onClick={() => setDepth(0)}><Icon>chevron_left</Icon>Back</MenuItem>
                    {accountStore.accounts.filter((account) => !account.companyId)
                        .map((account) => <MenuAccount key={account.id} {...account}/>)}
                    {accountStore.companies.filter((company) => company.accounts.length > 0)
                        .map((company) => <MenuCompany key={company.id} company={company} callback={selectCompany} />)}
                </MenuList>
                <MenuList ref={(el) => setWidth('accounts', el?.clientWidth ?? 0)}>
                    <MenuItem onClick={() => setDepth(1)}><Icon>chevron_left</Icon>Back</MenuItem>
                    {company ? company.accounts.map((account) => <MenuAccount key={account.id} {...account} />) : null}
                </MenuList>
            </div>
        </div>);
}, {forwardRef: true});

export default PageMenu;
