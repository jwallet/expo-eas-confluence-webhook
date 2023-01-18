// https://blog.expo.dev/introducing-support-for-postpublish-hooks-46ab3258fb03

const axios = require('axios');

const TOKEN = process.env.TOKE; // 'Seek3LESkZVxsbSiKC8n1C2B';
const USER_NAME = process.env.USER; // 'jose.ouellet@vooban.com';
const CONFLUENCE_CLOUD_DOMAIN = 'vooban'

type Page = {
  id: string;
  type: 'page';
  status: 'current';
  title: string;
  version: { number: number };
  body: { storage: { value: string; representation: 'storage' } };
};

const getConfluencePage = async (id: number): Promise<Page> => {
  const response = await axios.get(`https://${CONFLUENCE_CLOUD_DOMAIN}.atlassian.net/wiki/rest/api/content/${id}`, {
    params: {
      expand: 'body.storage,version',
    },
    auth: {
      username: USER_NAME,
      password: TOKEN,
    },
  });
  return response.data;
};

const updateConfluencePage = async (id: number, content: object) => {
  const response = await axios.put(`https://${CONFLUENCE_CLOUD_DOMAIN}.atlassian.net/wiki/rest/api/content/${id}`, content, {
    auth: {
      username: USER_NAME,
      password: TOKEN,
    },
  });
  return response.data;
};

type Build = {
  type?: 'Android' | 'iOS';
  id?: string;
  version?: string;
  sdk?: string;
  completedAt?: string;
  expiresAt?: string;
};

type Environment = 'continuous' | 'integration' | 'staging' | 'review';

const builds: Record<Environment, Build> = {
  continuous: {},
  integration: {},
  staging: {},
  review: {},
};

const generateTemplate = (build: Build) => {
  const buildURL = `https://expo.dev/accounts/guay/projects/guay/builds/${build.id}`;
  return `<table data-layout=\"default\" ac:local-id=\"da27978f-8ee4-4349-bf9b-9dc2fcb41544\"><colgroup><col style=\"width: 372.0px;\" /><col style=\"width: 198.0px;\" /><col style=\"width: 190.0px;\" /></colgroup><tbody><tr><td><p><strong>${build.type}</strong></p></td><td><p style=\"text-align: center;\"><strong>${build.version}</strong></p></td><td><p style=\"text-align: center;\"><strong>SDK ${build.sdk}</strong></p></td></tr><tr><td colspan=\"3\"><ac:structured-macro ac:name=\"iframe\" ac:schema-version=\"1\" data-layout=\"default\" ac:local-id=\"a8b49c3e-5066-4b7f-b8b9-408f14c246bc\" ac:macro-id=\"96338a9b108fa6f425a11b745bfac3f9\"><ac:parameter ac:name=\"scrolling\">no</ac:parameter><ac:parameter ac:name=\"src\"><ri:url ri:value=\"https://api.qrserver.com/v1/create-qr-code/?size=200x200&amp;data=${buildURL}\" /></ac:parameter><ac:parameter ac:name=\"width\">200</ac:parameter><ac:parameter ac:name=\"frameborder\">hide</ac:parameter><ac:parameter ac:name=\"align\">middle</ac:parameter><ac:parameter ac:name=\"title\">QR Code</ac:parameter><ac:parameter ac:name=\"longdesc\">Scan QR Code to install</ac:parameter><ac:parameter ac:name=\"height\">200</ac:parameter></ac:structured-macro></td></tr><tr><td colspan=\"3\"><p><a href=\"${buildURL}" data-card-appearance=\"inline\">${buildURL}</a> </p></td></tr><tr><td><p>Completed at: <strong>${build.completedAt}/strong></p></td><td colspan=\"2\"><p style=\"text-align: right;\">Expirates at: <strong>${build.expiresAt}</strong></p></td></tr></tbody></table><p />`;
};

type ExpoBuild = {
  id: string;
  platform: 'ios' | 'android';
  status: 'finished' | 'errored' | 'cancelled';
  buildDetailsPageUrl: string;
  artifacts: {
    buildUrl: string;
  };
  metadata: {
    appVersion: string;
    appBuildVersion: string;
    buildProfile: Environment;
    sdkVersion: string;
  };
  completedAt: string;
  expirationDate: string;
};

enum PagePerBuildProfile {
  ['continuous:android'] = 3317694616,
  ['continuous:ios'] = 3317694616,
  ['integration:android'] = 0,
  ['integration:ios'] = 0,
  ['staging:android'] = 0,
  ['staging:ios'] = 0,
  ['review:android'] = 0,
  ['review:ios'] = 0,
}

type ExpoPostPublishHookOptions = {
  url: string; // Published URL of the project
  iosBundle: string; // iOS JS bundle as a string
  iosSourceMap: string; // iOS source map as a string
  iosManifest: any; // Published iOS manifest
  androidBundle: string; // Android JS bundle as a string
  androidSourceMap: string; // Android source map as a string
  androidManifest: any; // Published Android manifest
  projectRoot: string; // Path to the project
};

const main = async (context: ExpoBuild) => {
  const pageId = PagePerBuildProfile[`${context.metadata.buildProfile}:${context.platform}`];
  const page = await getConfluencePage(pageId);
  await updateConfluencePage(pageId, {
    version: {
      number: page.version.number + 1,
      message: `EAS build ${context.id} finished`,
    },
    type: 'page',
    status: 'current',
    title: page.title,
    body: {
      storage: {
        value: generateTemplate({
          type: context.platform === 'ios' ? 'iOS' : 'Android',
          id: context.id,
          version: context.metadata.appVersion,
          sdk: context.metadata.sdkVersion,
          completedAt: context.completedAt,
          expiresAt: context.expirationDate,
        }),
        representation: 'storage',
      },
    },
  });
};

export default main;