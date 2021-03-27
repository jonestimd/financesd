import defaultTranslations from './defaults.json';

export const defaultKey = 'defaults';

export interface ITranslations {
    [key: string]: string;
}

export interface IBundle {
    [locale: string]: ITranslations;
}

const defaultBundle = {
    [defaultKey]: defaultTranslations,
};

export function translate(key: string, bundle: IBundle = defaultBundle) {
    const translations = bundle[navigator.language];
    if (translations && translations[key]) return translations[key];
    return bundle[defaultKey][key] ?? key;
}
