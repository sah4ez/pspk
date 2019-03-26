using System;
using System.Net.Http;
using System.Text;
using Android.App;
using Android.Widget;
using Android.OS;
using Android.Views;
using Newtonsoft.Json.Linq;
using String = Java.Lang.String;

namespace PSP_Key
{
    [Activity(Theme = "@android:style/Theme.Material.Light", Label = "PSP_Key", MainLauncher = true, Icon = "@drawable/icon")]
    public class MainActivity : Activity
    {
        private TextView textView;
        private LinearLayout _getLayout { get; set; }
        private EditText editTextName { get; set; }
        private EditText editTextKey { get; set; }
        private Button btnSend { get; set; }
        private Button btnGet { get; set; }

        protected override void OnCreate(Bundle bundle)
        {
            base.OnCreate(bundle);

            SetContentView(Resource.Layout.Main);

            _getLayout = FindViewById<LinearLayout>(Resource.Id.getLayout);
            editTextName = FindViewById<EditText>(Resource.Id.edittext_getByName);
            editTextKey = FindViewById<EditText>(Resource.Id.edittext_getByKey);
            btnSend = FindViewById<Button>(Resource.Id.buttonSet);
            btnGet = FindViewById<Button>(Resource.Id.buttonGet);

            var client = new HttpClient();
            var host = new Uri("https://pspk.now.sh/");
            
            btnSend.Click += (sender, args) =>
            {
                
                var name = editTextName.Text;
                var key = editTextKey.Text;
                var content = new StringContent($"{{ \"name\": \"{name}\", \"key\": \"{key}\" }}");

                var response = client.PostAsync(host, content).GetAwaiter().GetResult();
                response.Content.ReadAsStringAsync().GetAwaiter().GetResult();
            };
            btnGet.Click += (sender, args) =>
            {
                var name = editTextName.Text;
                var content = new StringContent($"{{ \"name\": \"{name}\" }}");

                var response = client.PostAsync(host, content).GetAwaiter().GetResult();
                var text = response.Content.ReadAsStringAsync().GetAwaiter().GetResult();
                try
                {
                    var json = JToken.Parse(text);
                    editTextKey.Text = json["key"].ToObject<string>();
                }
                catch (Exception ex)
                {
                }
            };
//            ActionBar.NavigationMode = ActionBarNavigationMode.Tabs;

//            AddTab("GET");  
//            AddTab("SET");
        }
        
        private void AddTab(string tabText)  
        {  
            ActionBar.Tab tab = ActionBar.NewTab();  
            tab.SetText(tabText);  
            tab.TabSelected += OnTabSelected;  
            ActionBar.AddTab(tab);  
        }
        
        private void OnTabSelected(object sender, ActionBar.TabEventArgs args)  
        {  
            var CurrentTab = (ActionBar.Tab)sender;

            if (CurrentTab.Position == 0)
            {
//                textView.Text = "Tab One Selected";
//                _setLayout.Visibility = ViewStates.Invisible;
//                _getLayout.Visibility = ViewStates.Visible;
            }

            else  
            {  
//                textView.Text = "Tab Two Selected";  
//                _setLayout.Visibility = ViewStates.Visible;
//                _getLayout.Visibility = ViewStates.Invisible;
            }  
        } 
    }
}