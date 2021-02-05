import React from 'react';
import {Link} from 'react-router-dom';
import {MenuList, MenuItem, ListItemText, IconButton, Icon} from '@material-ui/core';
import {translate} from '../../i18n/localize';

const menuItems = [
    ['menu.categories', '/finances/categories'],
    ['menu.groups', '/finances/groups'],
    ['menu.payees', '/finances/payees'],
    ['menu.securities', '/finances/securities'],
];

interface IProps {
    currentPage: string;
    setDepth: (depth: -1 | 1) => void;
}

const MainMenu: React.FC<IProps> = ({currentPage, setDepth}) => {
    const items = menuItems.filter(([key]) => currentPage !== key);
    return (
        <MenuList>
            {currentPage !== 'menu.accounts' ?
                <MenuItem>
                    <ListItemText><Link to='/finances/'>{translate('menu.accounts')}</Link></ListItemText>
                    <IconButton size='small' onClick={() => setDepth(1)}><Icon>chevron_right</Icon></IconButton>
                </MenuItem> : null}
            {items.map(([key, url]) => (<MenuItem key={key}><Link to={url}>{translate(key)}</Link></MenuItem>))}
        </MenuList>
    );
};

export default MainMenu;
