function Owb(){}
function Iwb(){}
function vRc(){}
function MRc(){}
function WRc(){}
function RRc(){}
function ORc(b){this.b=b}
function CRc(){DRc.call(this,_Rc())}
function DRc(){ERc.call(this,hsc((_Rc(),$Rc)))}
function _Rc(){_Rc=Xed;$Rc=new qFb(HVd,7,7)}
function XRc(b){b.style[kld]=lld;b.style[old]=Zsd;b.style[rld]=Zsd}
function VRc(b,c,d,e,f){b.style[wld]=e+Whd;c.style[pld]=e+Whd;d.style[pld]=f+Whd}
function BRc(b,c){var d;b.f=c;d=b.g[0];d.style[wld]=c;URc(b.c,parseInt(d[DFd])||0)}
function zRc(b,c){GQc(b.N,chd,c);GQc(b.j,c,DVd);GQc(b.g[0],c,pld);GQc(b.g[1],c,sld)}
function Nwb(){var b;while(Jwb){b=Jwb;Jwb=Jwb.c;!Jwb&&(Kwb=null);bHb(b.b.b,Y3b())}}
function Qwb(){Lwb=new Owb;TB((QB(),PB),33);!!$stats&&$stats(KC(FWd,mhd,-1,-1));Lwb.Cd();!!$stats&&$stats(KC(FWd,ODd,-1,-1))}
function URc(b,c){var d,e,f,g;f=b.b.j;e=parseInt(b.b.b[DFd])||0;g=parseInt(f[DFd])||0;if(e<g){return}d=e-c-g;c<0?(c=0):d<0&&(c=e-g);VRc(b.b.g[0],f,b.b.g[1],c,c+g)}
function Y3b(){var b,c,d,e;e=new CRc;zRc(e,GWd);e.N.style[vld]=mOd;e.N.style[wld]=AVd;BRc(e,BVd);d=HWd;for(c=0;c<2;++c){d+=d}wDc(e,0,new yxc(d));wDc(e,1,new yxc(d));b=new vwc;Vv(b,e);return b}
function ERc(b){var c,d,e,f,g;yDc.call(this,$doc.createElement($gd),$doc.createElement($gd),EDc($doc.createElement($gd)),EDc($doc.createElement($gd)));this.c=new WRc;this.b=EDc($doc.createElement($gd));e=this.g[0];c=this.g[1];d=this.j;this.N.appendChild(this.b);this.b.appendChild(e);this.b.appendChild(d);this.b.appendChild(c);d.innerHTML=IWd+QSc(b.e,b.c,b.d,b.f,b.b)+Hrd||chd;e.style[mld]=Psd;c.style[mld]=Psd;this.N[Rgd]=JWd;this.c.b=this;this.N.style[kld]=tld;g=this.g[0];f=this.g[1];XRc(g);XRc(f);XRc(this.j);zDc(this.b);f.style[sld]=Zsd;BRc(this,GVd)}
var IWd="<div class='vsplitter' style='text-align:center;'>",KWd='AsyncLoader33',MWd='VerticalSplitPanel',NWd='VerticalSplitPanel$1',LWd='VerticalSplitPanel$Impl',HWd='Voici un texte permettant de voir comment le contenu situ\xE9 de chaque c\xF4t\xE9 de la barre de fractionnement se d\xE9file. ',GWd='cwVerticalSplitPanel',JWd='gwt-VerticalSplitPanel',FWd='runCallbacks33';_=Owb.prototype=Iwb.prototype=new hr;_.gC=function Pwb(){return d8};_.Cd=function Twb(){Nwb()};_.cM={};_=CRc.prototype=vRc.prototype=new qDc;_.gC=function FRc(){return pib};_.ub=function GRc(b){zRc(this,b)};_.Gb=function HRc(){BRc(this,this.f);Hoc(new ORc(this))};_.of=function IRc(b,c){URc(this.c,this.e+c-this.d)};_.pf=function JRc(b,c){this.d=c;this.e=parseInt(this.g[0][DFd])||0};_.Hb=function KRc(){};_.vb=function LRc(b){this.N.style[wld]=b};_.cM={13:1,14:1,16:1,17:1,75:1,77:1,126:1};_.b=null;_.d=0;_.e=0;_.f=null;_=ORc.prototype=MRc.prototype=new hr;_.$b=function PRc(){BRc(this.b,this.b.f)};_.gC=function QRc(){return nib};_.cM={50:1,108:1};_.b=null;_=WRc.prototype=RRc.prototype=new hr;_.gC=function YRc(){return oib};_.cM={};_.b=null;var $Rc=null;var d8=U_c(owd,KWd),oib=U_c($td,LWd),pib=U_c($td,MWd),nib=U_c($td,NWd);$entry(Qwb)();