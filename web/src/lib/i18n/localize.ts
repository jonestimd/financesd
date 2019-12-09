import defaultTranslations from './defaults.json';

export const defaultKey = 'defaults';

export type Translations = {[key: string]: string}

export type Bundle = {[locale: string]: Translations};

const defaultBundle = {
    [defaultKey]: defaultTranslations
};

export function translate(key: string, bundle: Bundle = defaultBundle) {
    const translations = bundle[navigator.language];
    if (translations && translations[key]) return translations[key];
    return bundle[defaultKey][key] || key;
}