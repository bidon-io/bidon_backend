module RequestParams
  def config_params
    {
      device:   device_params,
      session:  session_params,
      app:      app_params,
      user:     user_params,
      geo:      geo_params,
      adapters: adapter_params,
      ext:      default_object_params,
      token:    '{}',
    }
  end

  def auction_params
    {
      device:     device_params,
      session:    session_params,
      app:        app_params,
      user:       user_params,
      geo:        geo_params,
      adapters:   adapter_params,
      ext:        default_object_params,
      ad_object:  default_object_params,
      segment_id: 'some segment id',
      token:      '{}',
    }
  end

  def stats_params
    {
      stats:      {
        auction_id:               'f26af577-869e-41cb-909e-4d3eba57a28b',
        auction_configuration_id: 10,
        rounds:                   [
          {
            id:          'postbid',
            pricefloor:  1.0,
            winner_id:   'bidmachine',
            winner_ecpm: 1.0,
            demands:     [
              {
                id:             'admob',
                ad_unit_id:     'AAAA',
                status:         'WIN',
                ecpm:           1.0,
                bid_start_ts:   123,
                bid_finish_ts:  124,
                fill_start_ts:  126,
                fill_finish_ts: 130,
              },
            ],
          },
        ],
      },
      device:     device_params,
      session:    session_params,
      app:        app_params,
      user:       user_params,
      geo:        geo_params,
      ext:        default_object_params,
      token:      '{}',
      segment_id: 'some segment id',
    }
  end

  def show_params
    {
      show:       {
        auction_id:               'f26af577-869e-41cb-909e-4d3eba57a28b',
        auction_configuration_id: 10,
        imp_id:                   '66b039f6-d43a-49ee-a84d-1eee15e91fba',
        demand_id:                'admob',
        ad_unit_id:               'AAAA',
        ecpm:                     1.0,
        banner:                   {
          format: 'LEADERBOARD',
        },
        interstitial:             {
        },
        rewarded:                 {
        },
      },
      device:     device_params,
      session:    session_params,
      app:        app_params,
      user:       user_params,
      geo:        geo_params,
      ext:        default_object_params,
      token:      '{}',
      segment_id: 'some segment id',
    }
  end

  def device_params
    {
      ua:              'User Agent',
      make:            'Apple',
      model:           'iPhone',
      os:              'iOS',
      osv:             '15.0.0',
      hwv:             '14,2',
      h:               2532,
      w:               1170,
      ppi:             2,
      pxratio:         3.0,
      js:              1,
      language:        'en',
      carrier:         'Orange',
      mccmnc:          '210-102',
      connection_type: 'WIFI',
    }
  end

  def session_params
    {
      id:                           '51acc730-1402-11ed-861d-0242ac120002',
      launch_ts:                    '1659571550',
      launch_monotonic_ts:          '1203445',
      start_ts:                     1_659_571_550,
      monotonic_start_ts:           1_203_445,
      ts:                           1_659_571_594,
      monotonic_ts:                 1_203_497,
      memory_warnings_ts:           [
        1_659_571_572,
      ],
      memory_warnings_monotonic_ts: [
        1_203_464,
      ],
      ram_used:                     102_858_752,
      ram_size:                     5_971_034_112,
      storage_free:                 16_699_088_896,
      storage_used:                 111_182_376_960,
      battery:                      0.86,
      cpu_usage:                    0.24,
    }
  end

  def app_params
    {
      bundle:            'myamazing.app.com',
      key:               'some key',
      framework:         'unity',
      version:           '1.2.3',
      framework_version: '14.3.2',
      plugin_version:    '1.2.3',
    }
  end

  def user_params
    {
      idfa:                          'UUID',
      tracking_authorization_status: 3,
      idfv:                          'UUID',
      idg:                           'UUID',
      consent:                       {
        key1: 'value1',
      },
      coppa:                         false,
    }
  end

  def geo_params
    {
      lat:       23.12,
      lon:       -45.95,
      accuracy:  10,
      lastfix:   23,
      country:   'PL',
      city:      'Warsaw',
      zip:       '02-235',
      utcoffset: -432_000,
    }
  end

  def adapter_params
    {
      admob:      {
        version:     '0.1.0.2',
        sdk_version: '7.9.0',
      },
      bidmachine: {
        version:     '0.1.0.2',
        sdk_version: '7.9.0',
      },
      applovin:   {
        version:     '0.1.0.2',
        sdk_version: '7.9.0',
      },
      appsflyer:  {
        version:     '0.1.0.2',
        sdk_version: '7.9.0',
      },
    }
  end

  def default_object_params
    {
      key1: 'value1',
    }
  end
end
